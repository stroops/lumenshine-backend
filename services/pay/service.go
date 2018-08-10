package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/environment"
	"math/big"
	"net"

	"github.com/Soneso/lumenshine-backend/services/pay/bitcoin"
	"github.com/Soneso/lumenshine-backend/services/pay/ethereum"
	"github.com/Soneso/lumenshine-backend/services/pay/stellar"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//our server
type server struct {
	Env *environment.Environment
}

//StartGrpcService will start the gRPC Service
func StartGrpcService(env *environment.Environment, log *logrus.Entry) error {
	//start the service

	cnf := env.Config
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"port": cnf.Port}).Fatalf("Failed to listen")
	}
	log.WithFields(logrus.Fields{"port": cnf.Port}).Print("Pay-Service listening")

	//v, err := stellar.XLMToStroops(0.0000001)
	/*v, err := ethereum.EthToWei(0.00000001)
	if err != nil {
		panic(err)
	}
	log.Info("v:", v)*/

	s := grpc.NewServer()
	pb.RegisterPayServiceServer(s, &server{
		Env: env,
	})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatalf("Failed to serve")
	}

	return nil
}

func (s *server) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	var err error
	chainAmount, chainAmountDenom, err := s.getPriceForCoins(r.Chain, r.CoinAmount)
	if err != nil {
		log.WithError(err).Error("Error getting price")
		return nil, err
	}

	if chainAmount == 0.0 {
		return nil, errors.New("Could not get Chain amount")
	}
	ret := &pb.CreateOrderResponse{
		ChainAmount: chainAmount,
		Chain:       r.Chain,
	}

	addressSeed := ""
	var index uint32
	index = 0
	if r.Chain == m.ChainEth || r.Chain == m.ChainBTC {
		//get new chain address
		index, err = env.DBC.GetNextChainAddressIndex(r.Chain)
		if r.Chain == m.ChainEth {
			ret.Address, err = s.Env.EthereumAddressGenerator.Generate(index)
		} else {
			ret.Address, err = s.Env.BitcoinAddressGenerator.Generate(index)
		}
	} else if r.Chain == m.ChainXLM {
		ret.Address, addressSeed, err = s.Env.StellarAddressGenerator.Generate()
	} else if r.Chain == m.ChainFiat {
		//return the fiat data
		ret.FiatBic = s.Env.Config.Fiat.BIC
		ret.FiatIban = s.Env.Config.Fiat.IBAN
		ret.FiatDestinationName = s.Env.Config.Fiat.DestiantionName
		ret.FiatPaymentUsage = s.Env.Config.Fiat.PaymentUsage
	} else {
		return nil, errors.New("Wrong chain")
	}

	if err != nil {
		return nil, err
	}

	//save order to db
	order := &m.UserOrder{
		UserID:               int(r.UserId),
		OrderStatus:          m.OrderStatusWaitingForPayment,
		CoinAmount:           r.CoinAmount,
		ChainAmount:          chainAmount,
		ChainAmountDenom:     chainAmountDenom.String(),
		Chain:                r.Chain,
		AddressIndex:         int64(index),
		ChainAddress:         ret.Address,
		ChainAddressSeed:     addressSeed,
		UserStellarPublicKey: r.UserPublicKey,
		OrderPhaseID:         r.IcoPhase,
	}
	err = order.Insert(s.Env.DBC)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain": r.Chain, "address": ret.Address, "index": index,
		}).Error("Could not create order")
	}

	ret.OrderId = int64(order.ID)
	return ret, err
}

// getPriceForCoins returns the price for the given coin amount in the corresponding chain value
func (s *server) getPriceForCoins(chain string, coinCount int64) (chainAmount float64, chainDenomAmount *big.Int, err error) {
	chainAmount = 0.0
	chainDenomAmount = big.NewInt(0)

	if chain == m.ChainEth {
		chainAmount = s.Env.Config.Ethereum.TokenPrice * float64(coinCount)
		chainDenomAmount, err = ethereum.EthToWei(fmt.Sprintf("%f", chainAmount))
	} else if chain == m.ChainBTC {
		chainAmount = s.Env.Config.Bitcoin.TokenPrice * float64(coinCount)
		chainDenomAmount, err = bitcoin.BtcToSat(fmt.Sprintf("%f", chainAmount))
	} else if chain == m.ChainXLM {
		chainAmount = s.Env.Config.Stellar.TokenPrice * float64(coinCount)
		chainDenomAmount, err = stellar.XLMToStroops(fmt.Sprintf("%f", chainAmount))
	} else if chain == m.ChainFiat {
		chainAmount = s.Env.Config.Fiat.TokenPrice * float64(coinCount)
	}

	return
}

func (s *server) GetCoinPrice(ctx context.Context, r *pb.CoinPriceRequest) (*pb.CoinPriceResponse, error) {
	chainAmount, chainDenomAmount, err := s.getPriceForCoins(r.Chain, r.CoinAmount)
	return &pb.CoinPriceResponse{
		ChainAmount:             chainAmount,
		ChainAmountDenomination: chainDenomAmount.String(),
	}, err
}

func (s *server) GetUserOrders(ctx context.Context, r *pb.UserOrdersRequest) (*pb.UserOrdersResponse, error) {
	q := []qm.QueryMod{
		qm.Where(m.UserOrderColumns.UserID+"=?", r.UserId),
	}

	if r.Status != "" {
		q = append(q, qm.Where(m.UserOrderColumns.OrderStatus+"=?", r.Status))
	}

	if r.OrderId != 0 {
		q = append(q, qm.Where("id=?", r.OrderId))
	}

	orders, err := m.UserOrders(s.Env.DBC, q...).All()
	if err != nil {
		return nil, err
	}
	ret := new(pb.UserOrdersResponse)
	ret.UserOrders = make([]*pb.UserOrder, len(orders))
	for i := 0; i < len(orders); i++ {
		ret.UserOrders[i] = &pb.UserOrder{
			Id:                   int64(orders[i].ID),
			OrderStatus:          orders[i].OrderStatus,
			CoinAmount:           orders[i].CoinAmount,
			ChainAmount:          orders[i].ChainAmount,
			ChainAmountDenom:     orders[i].ChainAmountDenom,
			Chain:                orders[i].Chain,
			ChainAddress:         orders[i].ChainAddress,
			UserStellarPublicKey: orders[i].UserStellarPublicKey,
		}
		if orders[i].Chain == m.ChainFiat {
			ret.UserOrders[i].FiatBic = s.Env.Config.Fiat.BIC
			ret.UserOrders[i].FiatIban = s.Env.Config.Fiat.IBAN
			ret.UserOrders[i].FiatDestinationName = s.Env.Config.Fiat.DestiantionName
			ret.UserOrders[i].FiatPaymentUsage = s.Env.Config.Fiat.PaymentUsage
		}
	}

	return ret, nil
}

func (s *server) UserOrderSetStatus(ctx context.Context, r *pb.UserOrderSetStatusRequest) (*pb.Empty, error) {
	order, err := m.UserOrders(s.Env.DBC, qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One()
	if err != nil {
		return nil, err
	}

	order.OrderStatus = r.Status
	order.PaymentErrorMessage = r.ErrorMessage
	order.UpdatedBy = r.Base.UpdateBy

	err = order.Update(s.Env.DBC,
		m.UserOrderColumns.OrderStatus,
		m.UserOrderColumns.PaymentErrorMessage,
		m.UserOrderColumns.UpdatedAt,
		m.UserOrderColumns.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) GetActveICOPhase(ctx context.Context, r *pb.Empty) (*pb.IcoPhaseResponse, error) {
	p, err := m.IcoPhases(s.Env.DBC, qm.Where(m.IcoPhaseColumns.IsActive+"=?", true)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.IcoPhaseResponse{}, nil
		}
		return nil, err
	}

	return &pb.IcoPhaseResponse{
		PhaseName:  p.PhaseName,
		StartTime:  int64(p.StartTime.Unix()),
		EndTime:    int64(p.EndTime.Unix()),
		CoinAmount: p.CoinAmount,
		IsActive:   p.IsActive,
	}, nil
}

func (s *server) PayGetTrustStatus(ctx context.Context, r *pb.PayGetTrustStatusRequest) (*pb.BoolResponse, error) {
	order, err := m.UserOrders(s.Env.DBC, qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One()
	if err != nil {
		return nil, err
	}

	hasTrust, err := s.Env.AccountConfigurator.GetTrustStatus(order)
	return &pb.BoolResponse{Value: hasTrust}, err
}

func (s *server) PayGetTransaction(ctx context.Context, r *pb.PayGetTransactionRequest) (*pb.PayGetTransactionResponse, error) {
	order, err := m.UserOrders(s.Env.DBC, qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One()
	if err != nil {
		return nil, err
	}

	tx, errCode, err := s.Env.AccountConfigurator.GetPaymentTransaction(order)
	if err == nil && errCode == 0 {
		//update the order to hold the tx
		order.PaymentTX = tx
		order.UpdatedBy = r.Base.UpdateBy
		err = order.Update(s.Env.DBC,
			m.UserOrderColumns.PaymentTX,
			m.UserOrderColumns.UpdatedAt,
			m.UserOrderColumns.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
	}

	return &pb.PayGetTransactionResponse{
		ErrorCode:   errCode,
		Transaction: tx,
	}, err
}

func (s *server) PayExecuteTransaction(ctx context.Context, r *pb.PayExecuteTransactionRequest) (*pb.Empty, error) {
	order, err := m.UserOrders(s.Env.DBC, qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One()
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, s.Env.AccountConfigurator.ExecuteTransaction(order, r.Transaction)
}
