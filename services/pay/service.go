package main

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/db/modext"
	"github.com/Soneso/lumenshine-backend/services/pay/environment"
	context "golang.org/x/net/context"

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
			EcAssetIssuerPk:      ec.EcAssetIssuerPK,
			DenomPricePerToken:   denom.String(),
			PricePerToken:        nativAmount,
		})
	}

	return &pb.PhaseDataResponse{
		Id:                       int64(p.ID),
		IcoId:                    int64(p.IcoID),
		IcoPhaseName:             p.IcoPhaseName,
		IcoPhaseStatus:           p.IcoPhaseStatus,
		IcoIssuerPk:              p.R.Ico.IssuerPK,
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
		EcAssetIssuerPk:      p.EcAssetIssuerPK,
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
			// this is the public key for the payment recipient
			// for stellar this is one special account, where the user must withdraw his payment
			// this pk will be used in the order as the payment account
			paymentAddress = aec.StellarPaymentAccountPK
			paymentSeed = aec.StellarPaymentAccountSeed
		} else {
			// for other cryptos, we will generate a dedicated address for every order
			paymentAddress, paymentSeed, err = s.GeneratePaymentAddress(ec.PaymentNetwork, aec.ExchangeMasterKey)
			if err != nil {
				return nil, err
			}
		}
	}

	//save order to db
	//StellarUserPublicKey might be null. we will determin the account, when the payment arrives, in order to connect only once to horizen
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
	}
	err = o.Insert(s.Env.DBC, boil.Infer())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"payment_network": ec.PaymentNetwork, "address": paymentAddress, "index": addressIndex,
		}).Error("Could not create order")
		return nil, err
	}

	if ec.PaymentNetwork == m.PaymentNetworkStellar {
		o.PaymentUsage = fmt.Sprintf("%d", o.ID)
		_, err := o.Update(s.Env.DBC, boil.Whitelist(m.UserOrderColumns.PaymentUsage))
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"order-id": o.ID,
			}).Error("Could not update stellar payment usage")
			return nil, err
		}
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
		o.PaymentUsage = modext.UserOrderFiatPaymentUsage(aec.R.IcoPhaseBankAccount.PaymendUsageString, o)
		_, err := o.Update(s.Env.DBC, boil.Whitelist(m.UserOrderColumns.PaymentUsage))
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"order-id": o.ID,
			}).Error("Could not update payment usage")
			return nil, err
		}

		ret.FiatBic = aec.R.IcoPhaseBankAccount.BicSwift
		ret.FiatIban = aec.R.IcoPhaseBankAccount.Iban
		ret.FiatRecepientName = aec.R.IcoPhaseBankAccount.RecepientName
		ret.FiatBankName = aec.R.IcoPhaseBankAccount.BankName
	}

	ret.PaymentUsage = o.PaymentUsage

	if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeCrypto {
		//TODO
		//ret.PaymentQrImage = o.PaymentQRImage
		ret.PaymentAddress = paymentAddress
	}

	return ret, nil
}

func (s *server) GetUserOrders(ctx context.Context, r *pb.UserOrdersRequest) (*pb.UserOrdersResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	q := []qm.QueryMod{
		qm.Where(m.UserOrderColumns.UserID+"=?", r.UserId),
		//TODO: qm.Load(m.UserOrderRels.ProcessedTransaction),
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
			PaymentUsage:                       o.PaymentUsage,
		}

		/* TODO:
		if o.R.ProcessedTransaction != nil {
			if o.R.ProcessedTransaction.PaymentNetworkAmountDenomination != "" {
				denomReceived, err := ec.DenomFromString(o.R.ProcessedTransaction.PaymentNetworkAmountDenomination)
				if err != nil {
					return nil, err
				}
				ret.UserOrders[i].AmountReceived = ec.ToNativ(denomReceived)
			}
		}
		*/

		if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeFiat {
			ret.UserOrders[i].FiatBic = aec.R.IcoPhaseBankAccount.BicSwift
			ret.UserOrders[i].FiatIban = aec.R.IcoPhaseBankAccount.Iban
			ret.UserOrders[i].FiatRecepientName = aec.R.IcoPhaseBankAccount.RecepientName
			ret.UserOrders[i].FiatBankName = aec.R.IcoPhaseBankAccount.BankName
		}

		if ec.ExchangeCurrencyType == m.ExchangeCurrencyTypeCrypto {
			ret.UserOrders[i].PaymentAddress = o.PaymentAddress
			ret.UserOrders[i].StellarTransactionId = o.StellarTransactionID
			/* TODO:
			if o.R.ProcessedTransaction != nil {
				ret.UserOrders[i].PaymentRefundTxId = o.R.ProcessedTransaction.RefundTXID
			}
			*/

			//TODO
			//ret.UserOrders[i].PaymentQrImage = o.PaymentQRImage
		}
	}

	return ret, nil
}

func (s *server) PayGetTransaction(ctx context.Context, r *pb.PayGetTransactionRequest) (*pb.PayGetTransactionResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	db := s.Env.DBC
	var err error
	order, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One(db)
	if err != nil {
		return nil, err
	}

	ac := s.Env.AccountConfigurator

	if order.OrderStatus != m.OrderStatusWaitingUserTransaction {
		return &pb.PayGetTransactionResponse{ErrorCode: cerr.OrderWrongStatus}, nil
	}

	pk := ""
	var pkExists bool

	if order.StellarUserPublicKey != "" {
		_, pkExists, err = ac.GetAccount(order.StellarUserPublicKey)
		if err != nil {
			return nil, err
		}
		if !pkExists {
			return &pb.PayGetTransactionResponse{
				ErrorCode:        cerr.UserNotExists,
				StellarPublicKey: order.StellarUserPublicKey,
			}, nil
		}
		pk = order.StellarUserPublicKey
	} else {
		//no pk from order -> check wallets
		//we take the first funded wallet
		wallets, err := m.UserWallets(qm.Where(m.UserWalletColumns.UserID+"=?", r.UserId)).All(db)
		if err != nil {
			return nil, err
		}
		for _, w := range wallets {
			_, pkExists, err = ac.GetAccount(w.PublicKey)
			if err != nil {
				return nil, err
			}
			if pkExists {
				//we use this pk for the order
				pk = w.PublicKey
				break
			}
		}
		if pk == "" {
			//there is no existing pk in stellar. therefore we take the first one from the wallets
			pk = wallets[0].PublicKey

			//if no pk was found, we need to check, if we ever created one account for the user. if so, the user merged the account and we will return an error
			//we create only one account per user and lifetime
			u, err := m.UserProfiles(qm.Where("id=?", r.UserId)).One(db)
			if err != nil {
				return nil, err
			}
			if u.StellarAccountCreated {
				return &pb.PayGetTransactionResponse{
					ErrorCode:        cerr.UserShouldExist,
					StellarPublicKey: pk,
				}, nil
			}
		}
	}

	if !pkExists {
		//create the account and update the user profile, to reflect the stellar account creation
		err := ac.CreateAccount(pk, order)
		if err != nil {
			return nil, err
		}
	}

	//need to update the order to reflect the selected pk
	if order.StellarUserPublicKey == "" {
		order.StellarUserPublicKey = pk
		_, err := order.Update(db, boil.Whitelist(m.UserOrderColumns.StellarUserPublicKey))
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"order-id": order.ID,
			}).Error("Could not update stellar public key")
			return nil, err
		}
	}

	tx, errCode, err := s.Env.AccountConfigurator.GetPaymentTransaction(order)
	if err != nil {
		return nil, err
	}

	return &pb.PayGetTransactionResponse{
		ErrorCode:        errCode,
		Transaction:      tx,
		StellarPublicKey: order.StellarUserPublicKey,
	}, err
}

func (s *server) PayExecuteTransaction(ctx context.Context, r *pb.PayExecuteTransactionRequest) (*pb.PayExecuteTransactionResponse, error) {
	order, err := m.UserOrders(qm.Where(m.UserOrderColumns.UserID+"=? and id=?", r.UserId, r.OrderId)).One(s.Env.DBC)
	if err != nil {
		return nil, err
	}

	if order.OrderStatus != m.OrderStatusWaitingUserTransaction {
		return &pb.PayExecuteTransactionResponse{ErrorCode: cerr.OrderWrongStatus}, nil
	}

	hash, err := s.Env.AccountConfigurator.ExecuteTransaction(order, r.Transaction)
	return &pb.PayExecuteTransactionResponse{TxHash: hash}, err
}

func (s *server) FakePaymentTransaction(ctx context.Context, r *pb.TestTransaction) (*pb.BoolResponse, error) {
	log := helpers.GetDefaultLog(ServiceName, r.Base.RequestId)
	if s.Env.Config.AllowFakeTransactions {

		ch, ok := s.Env.Clients[r.PaymentChannel]
		if !ok {
			return nil, fmt.Errorf("PaymentChannel %s does not exist", r.PaymentChannel)
		}

		// Need to read the order first, because this is the normal "procedure". This will also update the order-status
		o, err := s.Env.DBC.GetOrderForAddress(ch, r.RecipientAddress, r.PaymentUsage)
		if err != nil {
			return &pb.BoolResponse{Value: false}, nil
		}

		if o == nil {
			return &pb.BoolResponse{Value: false}, errors.New("No order found")
		}

		tx := helpers.RandomString(10)
		if r.TxHash != "" {
			tx = r.TxHash
		}
		v := big.NewInt(r.DenomAmount)
		fmt.Println(v.String())

		ok, err = s.Env.DBC.AddNewTransaction(log, ch, tx, r.RecipientAddress, r.SenderAddress, o, v, 0)
		return &pb.BoolResponse{Value: ok}, err
	}

	return &pb.BoolResponse{Value: false}, nil
}
