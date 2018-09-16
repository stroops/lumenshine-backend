package main

import (
	"context"
	"fmt"
	"net"

	"github.com/Soneso/lumenshine-backend/helpers"
	"github.com/Soneso/lumenshine-backend/pb"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/db/modext"
	"github.com/Soneso/lumenshine-backend/services/pay/environment"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
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

func (s *server) GetPhaseData(ctx context.Context, r *pb.IDRequest) (*pb.PhaseDataResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	p, err := s.Env.DBC.GetICOPhaseByID(int(r.Id), log)
	if err != nil {
		return nil, err
	}

	ecs := []*pb.ExchangeCurrency{}

	for _, aec := range p.R.IcoPhaseActivatedExchangeCurrencies {
		ec, err := s.Env.DBC.GetExchangeCurrecnyByID(aec.R.ExchangeCurrency.ID, log)

		denom, nativAmount, err := s.Env.DBC.PriceForCoins(1, ec, p)
		if err != nil {
			return nil, err
		}

		ecs = append(ecs, &pb.ExchangeCurrency{
			Id:                   int64(ec.ID),
			ExchangeCurrencyType: ec.ExchangeCurrencyType,
			AssetCode:            ec.AssetCode,
			DenomAssetCode:       ec.DenomAssetCode,
			Decimals:             int64(ec.Decimals),
			PaymentNetwork:       ec.PaymentNetwork,
			IssuerPk:             ec.IssuerPK,
			DenomPricePerToken:   denom.String(),
			PricePerToken:        nativAmount,
		})
	}

	return &pb.PhaseDataResponse{
		Id:                       int64(p.ID),
		IcoId:                    int64(p.IcoID),
		IcoPhaseName:             p.IcoPhaseName,
		IcoPhaseStatus:           p.IcoPhaseStatus,
		StartTime:                p.StartTime.Unix(),
		EndTime:                  p.EndTime.Unix(),
		TokensLeft:               p.TokensLeft,
		TokenMaxOrderAmount:      p.TokenMaxOrderAmount,
		TokenMinOrderAmount:      p.TokenMinOrderAmount,
		ActiveExchangeCurrencies: ecs,
		MaxUserOrders:            int64(p.MaxUserOrders),
		IcoTokenAsset:            p.R.Ico.AssetCode,
	}, nil
}

func (s *server) GetUserOrderCount(ctx context.Context, r *pb.UserOrdersCountRequest) (*pb.IntResponse, error) {
	cnt, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and "+m.UserOrderColumns.IcoPhaseID+"=?", r.UserId, r.PhaseId)).Count(s.Env.DBC)
	if err != nil {
		return nil, err
	}
	return &pb.IntResponse{
		Value: cnt,
	}, nil
}

func (s *server) GetExchangeCurrencyData(ctx context.Context, r *pb.IDRequest) (*pb.ExchangeCurrency, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	p, err := s.Env.DBC.GetExchangeCurrecnyByID(int(r.Id), log)
	if err != nil {
		return nil, err
	}
	return &pb.ExchangeCurrency{
		Id:                   int64(p.ID),
		ExchangeCurrencyType: p.ExchangeCurrencyType,
		AssetCode:            p.AssetCode,
		DenomAssetCode:       p.DenomAssetCode,
		Decimals:             int64(p.Decimals),
		PaymentNetwork:       p.PaymentNetwork,
		IssuerPk:             p.IssuerPK,
	}, nil
}

func (s *server) GetCoinPrice(ctx context.Context, r *pb.CoinPriceRequest) (*pb.CoinPriceResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	ec, err := s.Env.DBC.GetExchangeCurrecnyByID(int(r.ExchangeCurrencyId), log)
	if err != nil {
		return nil, err
	}

	phase, err := s.Env.DBC.GetICOPhaseByID(int(r.IcoPhaseId), log)
	if err != nil {
		return nil, err
	}

	denomAmount, nativAmount, err := s.Env.DBC.PriceForCoins(r.CoinAmount, ec, phase)
	if err != nil {
		return nil, err
	}

	return &pb.CoinPriceResponse{
		ExchangeAmount:             nativAmount,
		ExchangeAmountDenomination: denomAmount.String(),
		ExchangeAssetCode:          ec.AssetCode,
	}, nil
}

func (s *server) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)

	phase, err := s.Env.DBC.GetICOPhaseByID(int(r.IcoPhaseId), log)
	if err != nil {
		return nil, err
	}

	ec, aec, err := s.Env.DBC.GetActiveExchangeCurrecnyByID(int(r.ExchangeCurrencyId), int(phase.ID), log)
	if err != nil {
		return nil, err
	}

	denomAmount, nativAmount, err := s.Env.DBC.PriceForCoins(r.TokenAmount, ec, phase)
	if err != nil {
		return nil, err
	}

	paymentAddress := ""
	paymentSeed := ""
	var addressIndex uint32

	if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeCrypto {
		if ec.PaymentNetwork == m.PaymentNetworkStellar {
			paymentAddress = aec.StellarPaymentAccountPK
			paymentSeed = aec.StellarPaymentAccountSeed
		} else {
			paymentAddress, paymentSeed, addressIndex, err = s.GeneratePaymentAddress(ec.PaymentNetwork, aec.ExchangeMasterKey)
			if err != nil {
				return nil, err
			}
		}

	}

	//save order to db
	o := &m.UserOrder{
		UserID:                             int(r.UserId),
		IcoPhaseID:                         phase.ID,
		OrderStatus:                        m.OrderStatusWaitingForPayment,
		TokenAmount:                        r.TokenAmount,
		StellarUserPublicKey:               r.UserPublicKey,
		ExchangeCurrencyID:                 ec.ID,
		ExchangeCurrencyDenominationAmount: denomAmount.String(),
		PaymentNetwork:                     ec.PaymentNetwork,
		PaymentAddress:                     paymentAddress,
		PaymentSeed:                        paymentSeed,
		AddressIndex:                       int64(addressIndex),
	}
	err = o.Insert(s.Env.DBC, boil.Infer())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"payment_network": ec.PaymentNetwork, "address": paymentAddress, "index": addressIndex,
		}).Error("Could not create order")
		return nil, err
	}

	ret := &pb.CreateOrderResponse{
		OrderId:                     int64(o.ID),
		OrderStatus:                 o.OrderStatus,
		ExchangeValueToPay:          nativAmount,
		ExchangeValueDenominator:    denomAmount.String(),
		ExchangeValueDenomAssetCode: ec.DenomAssetCode,
	}
	if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeFiat {
		//set the usage in the order
		o.FiatPaymentUsage = modext.UserOrderFiatPaymentUsage(aec.R.IcoPhaseBankAccount.PaymendUsageString, o)
		_, err := o.Update(s.Env.DBC, boil.Whitelist(m.UserOrderColumns.FiatPaymentUsage))
		if err != nil {
			//only log the error
			log.WithError(err).WithFields(logrus.Fields{
				"order-id": o.ID,
			}).Error("Could not update payment")
		}

		ret.FiatBic = aec.R.IcoPhaseBankAccount.BicSwift
		ret.FiatIban = aec.R.IcoPhaseBankAccount.Iban
		ret.FiatRecepientName = aec.R.IcoPhaseBankAccount.RecepientName
		ret.FiatPaymentUsage = o.FiatPaymentUsage
		ret.FiatBankName = aec.R.IcoPhaseBankAccount.BankName
	}

	if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeCrypto {
		ret.DepositPk = paymentAddress
		//TODO
		//ret.PaymentQrImage = o.PaymentQRImage
	}

	//TODO address generation
	return ret, nil
}

func (s *server) GetUserOrders(ctx context.Context, r *pb.UserOrdersRequest) (*pb.UserOrdersResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	q := []qm.QueryMod{
		qm.Where(m.UserOrderColumns.UserID+"=?", r.UserId),
	}

	if r.OrderStatus != "" {
		q = append(q, qm.Where(m.UserOrderColumns.OrderStatus+"=?", r.OrderStatus))
	}

	if r.OrderId != 0 {
		q = append(q, qm.Where("id=?", r.OrderId))
	}

	if r.IcoPhaseId != 0 {
		q = append(q, qm.Where(m.UserOrderColumns.IcoPhaseID+"=?", r.IcoPhaseId))
	}

	orders, err := m.UserOrders(q...).All(s.Env.DBC)
	if err != nil {
		return nil, err
	}
	ret := new(pb.UserOrdersResponse)
	ret.UserOrders = make([]*pb.UserOrder, len(orders))
	for i := 0; i < len(orders); i++ {
		o := orders[i]
		ec, aec, err := s.Env.DBC.GetActiveExchangeCurrecnyByID(o.ExchangeCurrencyID, o.IcoPhaseID, log)
		if err != nil {
			return nil, err
		}

		denom, err := ec.DenomFromString(o.ExchangeCurrencyDenominationAmount)
		if err != nil {
			return nil, err
		}

		ret.UserOrders[i] = &pb.UserOrder{
			Id:                                 int64(o.ID),
			OrderStatus:                        o.OrderStatus,
			IcoPhaseId:                         int64(o.IcoPhaseID),
			TokenAmount:                        o.TokenAmount,
			StellarUserPublicKey:               o.StellarUserPublicKey,
			ExchangeCurrencyId:                 int64(o.ExchangeCurrencyID),
			ExchangeCurrencyDenominationAmount: o.ExchangeCurrencyDenominationAmount,
			PaymentNetwork:                     o.PaymentNetwork,
			ExchangeAmount:                     ec.ToNativ(denom),
			ExchangeAssetCode:                  ec.AssetCode,
			ExchangeDenomAssetCode:             ec.DenomAssetCode,
			ExchangeCurrencyType:               ec.ExchangeCurrencyType,
		}

		if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeFiat {
			ret.UserOrders[i].FiatBic = aec.R.IcoPhaseBankAccount.BicSwift
			ret.UserOrders[i].FiatIban = aec.R.IcoPhaseBankAccount.Iban
			ret.UserOrders[i].FiatRecepientName = aec.R.IcoPhaseBankAccount.RecepientName
			ret.UserOrders[i].FiatPaymentUsage = o.FiatPaymentUsage
			ret.UserOrders[i].FiatBankName = aec.R.IcoPhaseBankAccount.BankName
		}

		if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeCrypto {
			ret.UserOrders[i].DepositPk = o.PaymentAddress
			ret.UserOrders[i].PaymentTxId = o.PaymentTXID
			ret.UserOrders[i].PaymentRefundTxId = o.PaymentRefundTXID
			//TODO
			//ret.UserOrders[i].PaymentQrImage = o.PaymentQRImage
		}
	}

	return ret, nil
}

func (s *server) PayGetTrustStatus(ctx context.Context, r *pb.PayGetTrustStatusRequest) (*pb.PayGetTrustStatusResponse, error) {
	order, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One(s.Env.DBC)
	if err != nil {
		return nil, err
	}

	hasTrust, err := s.Env.AccountConfigurator.GetTrustStatus(order)
	return &pb.PayGetTrustStatusResponse{
		HasStrust:            hasTrust,
		StellarAssetCode:     s.Env.Config.Stellar.TokenAssetCode,
		StellarIssuerAccount: s.Env.Config.Stellar.IssuerPublicKey,
	}, err
}

func (s *server) PayGetTransaction(ctx context.Context, r *pb.PayGetTransactionRequest) (*pb.PayGetTransactionResponse, error) {
	order, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One(s.Env.DBC)
	if err != nil {
		return nil, err
	}

	tx, errCode, err := s.Env.AccountConfigurator.GetPaymentTransaction(order)
	if err == nil && errCode == 0 {
		//update the order to hold the tx
		//TODO
		order.UpdatedBy = r.Base.UpdateBy
		_, err = order.Update(s.Env.DBC, boil.Whitelist(
			m.UserOrderColumns.UpdatedAt,
			m.UserOrderColumns.UpdatedBy,
		))
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
	order, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One(s.Env.DBC)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, s.Env.AccountConfigurator.ExecuteTransaction(order, r.Transaction)
}
