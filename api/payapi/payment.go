package main

import (
	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/pb"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/gin-gonic/gin"
)

// PriceForCoinRequest is the request data
// swagger:parameters PriceForCoinRequest PriceForCoin
type PriceForCoinRequest struct {
	// Amount of coins for the price-calculation
	// required: true
	CoinAmount int64 `json:"coin_amount" form:"coin_amount" query:"coin_amount" validate:"required,min=1"`

	// ID of the Exchange currency
	// required: true
	ExchangeCurrencyID int `json:"exchange_currency_id" form:"exchange_currency_id" query:"exchange_currency_id" validate:"required"`

	// ID of the ICO-Phase
	// required: true
	ICOPhaseID int `json:"ico_phase_id" form:"ico_phase_id" query:"ico_phase_id" validate:"required"`
}

// PriceForCoinResponse price for coin amount, based on the configuration
// swagger:model
type PriceForCoinResponse struct {
	ExchangeAmount    string `json:"exchange_amount"`
	ExchangeAssetCode string `json:"exchange_asset_code"`
}

// PriceForCoin returns the price for a given coin amount, related to the currency given
// swagger:route GET /ico_phase_price_for_amount PriceForCoin
//
// Returns the price for a given coin amount
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: PriceForCoinResponse
func PriceForCoin(uc *mw.IcopContext, c *gin.Context) {
	var l PriceForCoinRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	// correct state,amount, etc.
	price, err := payClient.GetCoinPrice(c, &pb.CoinPriceRequest{
		Base:               NewBaseRequest(uc),
		ExchangeCurrencyId: int64(l.ExchangeCurrencyID),
		CoinAmount:         l.CoinAmount,
		IcoPhaseId:         int64(l.ICOPhaseID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting price", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &PriceForCoinResponse{
		ExchangeAmount:    price.ExchangeAmount,
		ExchangeAssetCode: price.ExchangeAssetCode,
	})
}

// IcoPhaseDetailsRequest is the request data
// swagger:parameters IcoPhaseDetailsRequest IcoPhaseDetails
type IcoPhaseDetailsRequest struct {
	// ID of the ICO-Phase
	// required: true
	ICOPhaseID int `json:"ico_phase_id" form:"ico_phase_id" query:"ico_phase_id" validate:"required"`
}

// ExchangeCurrency represents one exchange currency used everywhere
// swagger:model
type ExchangeCurrency struct {
	ID                   int    `json:"id"`
	ExchangeCurrencyType string `json:"exchange_currency_type"`
	AssetCode            string `json:"asset_code"`

	// Number of decimal places for the denominator
	Decimals int64 `json:"decimals"`

	// Includes the UoM of the asset (EUR/XLM,BTC...)
	PricePerToken string `json:"price_per_token"`

	// this is the issuer public key for an stellar asset for the exchange currency. only set for stellar
	EcAssetIssuerPK string `json:"ec_asset_issuer_pk,omitempty"`
}

// IcoPhaseDetailsResponse lists details of the specified IPC-Phase
// swagger:model
type IcoPhaseDetailsResponse struct {
	ID                       int                `json:"id"`
	IcoID                    int                `json:"ico_id"`
	IcoPhaseName             string             `json:"ico_phase_name"`
	IcoPhaseStatus           string             `json:"ico_phase_status"`
	IcoIssuerPK              string             `json:"ico_issuer_pk"`
	StartTime                time.Time          `json:"start_time"`
	EndTime                  time.Time          `json:"end_time"`
	TokensLeft               int64              `json:"tokens_left"`
	TokenMaxOrderAmount      int64              `json:"token_max_order_amount"`
	TokenMinOrderAmount      int64              `json:"token_min_order_amount"`
	ActiveExchangeCurrencies []ExchangeCurrency `json:"active_exchange_currencies"`
}

// IcoPhaseDetails lists details of the specified IPC-Phase
// swagger:route GET /ico_phase_details IcoPhaseDetails
//
// Returns the details of a given ICO-Phase, including all activated Exchange-Currencies
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: IcoPhaseDetailsResponse
func IcoPhaseDetails(uc *mw.IcopContext, c *gin.Context) {
	var l IcoPhaseDetailsRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	// correct state,amount, etc.
	p, err := payClient.GetPhaseData(c, &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   int64(l.ICOPhaseID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting phase data", cerr.GeneralError))
		return
	}

	ret := &IcoPhaseDetailsResponse{
		ID:                  int(p.Id),
		IcoID:               int(p.IcoId),
		IcoPhaseName:        p.IcoPhaseName,
		IcoPhaseStatus:      p.IcoPhaseStatus,
		IcoIssuerPK:         p.IcoIssuerPk,
		StartTime:           time.Unix(p.StartTime, 0),
		EndTime:             time.Unix(p.EndTime, 0),
		TokensLeft:          p.TokensLeft,
		TokenMaxOrderAmount: p.TokenMaxOrderAmount,
		TokenMinOrderAmount: p.TokenMinOrderAmount,
	}
	ret.ActiveExchangeCurrencies = make([]ExchangeCurrency, len(p.ActiveExchangeCurrencies))
	for i, aec := range p.ActiveExchangeCurrencies {
		ret.ActiveExchangeCurrencies[i].AssetCode = aec.AssetCode
		ret.ActiveExchangeCurrencies[i].Decimals = aec.Decimals
		ret.ActiveExchangeCurrencies[i].ExchangeCurrencyType = aec.ExchangeCurrencyType
		ret.ActiveExchangeCurrencies[i].ID = int(aec.Id)
		ret.ActiveExchangeCurrencies[i].PricePerToken = aec.PricePerToken
		ret.ActiveExchangeCurrencies[i].EcAssetIssuerPK = aec.EcAssetIssuerPk
	}

	c.JSON(http.StatusOK, ret)
}

// CreateOrderRequest is the data for creating an order
// swagger:parameters CreateOrderRequest CreateOrder
type CreateOrderRequest struct {
	// ID of the Ico-Phase
	// required: true
	IcoPhaseID int64 `form:"ico_phase_id" json:"ico_phase_id" validate:"required"`

	// Ammount of tokens orderd
	// required: true
	OrderedTokenAmount int64 `form:"ordered_token_amount" json:"ordered_token_amount" validate:"required,min=1"`

	// ID of the Exchange currency
	// required: true
	ExchangeCurrencyID int64 `json:"exchange_currency_id" form:"exchange_currency_id" validate:"required"`

	// Stellar Public Key of the user for the payment/coins
	StellarUserPublicKey string `json:"stellar_user_public_key" form:"stellar_user_public_key"`
}

// CreateOrderResponse is the return data , for creating a new order
// swagger:model
type CreateOrderResponse struct {
	OrderID               int64  `json:"order_id"`
	OrderStatus           string `json:"order_status"`
	OrderedTokenAmount    int64  `json:"ordered_token_amount"`
	OrderedTokenAssetCode string `json:"ordered_token_asset_code"`
	PaymentNetwork        string `json:"payment_network"`

	//This is used for fiat and stellar payments. for stellar, this needs to be send via MEMO
	PaymentUsage string `json:"payment_usage,omitempty"`

	//StellarUserPublicKey is the stellar public key of the user for this order. If omited on CreateOrder, the service will grab the first 'free' one from the user wallets, when the payment arrives, in order to connect only once to horizen. So this might be empty
	StellarUserPublicKey string `json:"stellar_user_public_key"`

	//AssetCode in the payment Network
	ExchangeAssetCode string `json:"exchange_asset_code"`
	//Value to pay in the selected payment Network asset code
	ExchangeValueToPay string `json:"exchange_value_to_pay"`
	//Type for payment: stellar, other_crypto, fiat
	ExchangeCurrencyType string `json:"exchange_currency_type"`

	//PaymentAddress is the address in the PaymentNetwork, where the user must transfer the Exchange-Asset
	PaymentAddress string `json:"payment_address,omitempty"`

	//QRCode is a bitmap for the transaction in the Payment-Network
	//TODO
	QRCode []byte `json:"qr_code,omitempty"`

	FiatBIC             string `json:"fiat_bic,omitempty"`
	FiatIBAN            string `json:"fiat_iban,omitempty"`
	FiatDestinationName string `json:"fiat_destination_name,omitempty"`
	FiatBankName        string `json:"fiat_bank_name,omitempty"`
}

// CreateOrder creates a new order for the current user and returns the order-data for the next step
// swagger:route POST /create_order CreateOrder
//	   CreateOrder creates a new order for the current user and returns the order-data for the next step
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: CreateOrderResponse
func CreateOrder(uc *mw.IcopContext, c *gin.Context) {
	var l CreateOrderRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	phase, err := payClient.GetPhaseData(c, &pb.IDRequest{Base: NewBaseRequest(uc), Id: l.IcoPhaseID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting ico phase", cerr.GeneralError))
		return
	}

	if phase.IcoPhaseStatus != m.IcoStatusActive {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("Phase", cerr.NoActivePhase, "Phase is not active", ""))
		return
	}

	now := time.Now()
	pStart := time.Unix(phase.StartTime, 0)
	pEnd := time.Unix(phase.EndTime, 0)
	if now.Before(pStart) {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("Phase", cerr.NoActivePhase, "Phase not started yet", ""))
		return
	}

	if now.After(pEnd) {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("Phase", cerr.NoActivePhase, "Phase ended", ""))
		return
	}

	if l.OrderedTokenAmount < phase.TokenMinOrderAmount {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("TokenMinOrderAmount", cerr.OrderMinTokens, "To lees coins", ""))
		return
	}

	if l.OrderedTokenAmount > phase.TokenMaxOrderAmount {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("TokenMaxOrderAmount", cerr.OrderMaxTokens, "To much coins", ""))
		return
	}

	if l.OrderedTokenAmount > phase.TokensLeft {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("TokenMaxOrderAmount", cerr.InsufficientCoins, "Not enough coins left", ""))
		return
	}

	//check that passed exchange currency is active
	var ec *pb.ExchangeCurrency
	for _, aec := range phase.ActiveExchangeCurrencies {
		if aec.Id == l.ExchangeCurrencyID {
			ec = aec
		}
	}
	if ec == nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("exchange_currency_id", cerr.InvalidArgument, "ExchangeCurrency is not activated", ""))
		return
	}

	//check that user-order-count per phase is not excedeed
	cnt, err := payClient.GetUserOrderCount(c, &pb.UserOrdersCountRequest{
		UserId:  mw.GetAuthUser(c).UserID,
		PhaseId: l.IcoPhaseID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting user order count", cerr.GeneralError))
		return
	}

	if cnt.Value >= phase.MaxUserOrders {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("ordered_token_amount", cerr.ToMuchOrdersPerPhase, "To much orders for phase", ""))
		return
	}

	orderReq := &pb.CreateOrderRequest{
		Base:               NewBaseRequest(uc),
		UserId:             mw.GetAuthUser(c).UserID,
		IcoPhaseId:         l.IcoPhaseID,
		TokenAmount:        l.OrderedTokenAmount,
		ExchangeCurrencyId: l.ExchangeCurrencyID,
		UserPublicKey:      l.StellarUserPublicKey,
	}
	o, err := payClient.CreateOrder(c, orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating order", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &CreateOrderResponse{
		OrderID:               o.OrderId,
		OrderStatus:           o.OrderStatus,
		OrderedTokenAmount:    l.OrderedTokenAmount,
		OrderedTokenAssetCode: phase.IcoTokenAsset,
		PaymentNetwork:        ec.PaymentNetwork,
		PaymentUsage:          o.PaymentUsage,
		ExchangeAssetCode:     ec.AssetCode,
		ExchangeValueToPay:    o.ExchangeValueToPay,
		ExchangeCurrencyType:  ec.ExchangeCurrencyType,
		PaymentAddress:        o.PaymentAddress,
		QRCode:                o.PaymentQrImage,
		FiatBIC:               o.FiatBic,
		FiatIBAN:              o.FiatIban,
		FiatDestinationName:   o.FiatRecepientName,
		FiatBankName:          o.FiatBankName,
		StellarUserPublicKey:  o.UserPublicKey,
	})
}

// GetOrdersRequest is the data for filtering the user orders
// swagger:parameters GetOrdersRequest OrderList
type GetOrdersRequest struct {
	// ID of the Ico-Phase
	// required: false
	IcoPhaseID int64 `form:"ico_phase_id" json:"ico_phase_id"`

	// ID of the Order
	// required: false
	// OrderID int64 `form:"order_id" json:"order_id"`

	// ID of the Order
	// required: false
	OrderStatus string `form:"order_status" json:"order_status"`
}

// UserOrderResponse represents a UserOrder
// swagger:model
type UserOrderResponse struct {
	ID                   int64  `json:"id"`
	OrderStatus          string `json:"order_status"`
	IcoPhaseID           int64  `json:"ico_phase_id"`
	TokenAmount          int64  `json:"token_amount"`
	StellarUserPublicKey string `json:"stellar_user_public_key"`
	ExchangeCurrencyID   int64  `json:"exchange_currency_id"`

	//This is the exchange amount received in the first transaction
	AmountReceived string `json:"amount_received"`

	ExchangeCurrencyType string `json:"exchange_currency_type"`
	PaymentNetwork       string `json:"payment_network"`

	//used for fiat and stellar payments. for stellar, this must be send via MEMO
	PaymentUsage string `json:"payment_usage,omitempty"`

	//This is the public key in the PaymentNetwork, where the exchange-currency must be transfered to
	PaymentAddress string `json:"payment_address,omitempty"`

	//this is the coin payment tx in the stellar network
	StellarTransactionID string `json:"stellar_transaction_id,omitempty"`

	//this is the refund transaction id in the PaymentNetwork
	PaymentRefundTxID string `json:"payment_refund_tx_id,omitempty"`

	//TODO PaymentQrImage string `json:"payment_qr_image"`

	ExchangeAmount    string `json:"exchange_amount"`
	ExchangeAssetCode string `json:"exchange_asset_code"`

	FiatBic           string `json:"fiat_bic,omitempty"`
	FiatIban          string `json:"fiat_iban,omitempty"`
	FiatRecepientName string `json:"fiat_recepient_name,omitempty"`
	FiatBankName      string `json:"fiat_bank_name,omitempty"`
}

// OrderList returns the filtered list of orders for the current user.
// swagger:route GET /order_list OrderList
//     returns the filtered list of orders for the current user
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []UserOrderResponse
func OrderList(uc *mw.IcopContext, c *gin.Context) {
	var l GetOrdersRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	orders, err := payClient.GetUserOrders(c, &pb.UserOrdersRequest{
		Base:        NewBaseRequest(uc),
		UserId:      userID,
		OrderStatus: l.OrderStatus,
		IcoPhaseId:  l.IcoPhaseID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting orders", cerr.GeneralError))
		return
	}

	ret := make([]UserOrderResponse, len(orders.UserOrders))
	for i := 0; i < len(orders.UserOrders); i++ {
		ret[i] = getRespOrder(orders.UserOrders[i])
	}

	c.JSON(http.StatusOK, ret)
}

func getRespOrder(o *pb.UserOrder) UserOrderResponse {
	return UserOrderResponse{
		ID:                   o.Id,
		OrderStatus:          o.OrderStatus,
		TokenAmount:          o.TokenAmount,
		StellarUserPublicKey: o.StellarUserPublicKey,
		ExchangeCurrencyID:   o.ExchangeCurrencyId,
		ExchangeCurrencyType: o.ExchangeCurrencyType,
		PaymentNetwork:       o.PaymentNetwork,
		PaymentAddress:       o.PaymentAddress,
		StellarTransactionID: o.StellarTransactionId,
		PaymentRefundTxID:    o.PaymentRefundTxId,
		PaymentUsage:         o.PaymentUsage,

		//TODO PaymentQrImage : o.PaymentQrImage,

		ExchangeAmount:    o.ExchangeAmount,
		ExchangeAssetCode: o.ExchangeAssetCode,
		FiatBic:           o.FiatBic,
		FiatIban:          o.FiatIban,
		FiatRecepientName: o.FiatRecepientName,

		FiatBankName:   o.FiatBankName,
		AmountReceived: o.AmountReceived,
	}
}

// OrderDetailsRequest request-data
// swagger:parameters OrderDetailsRequest OrderDetails
type OrderDetailsRequest struct {
	// ID of the Order
	// required: true
	OrderID int64 `form:"order_id" json:"order_id" validate:"required"`
}

// OrderDetails returns the details for the specified order
// swagger:route GET /order_details OrderDetails
//     returns the details for the specified order
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: UserOrderResponse
func OrderDetails(uc *mw.IcopContext, c *gin.Context) {
	var l OrderDetailsRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID
	orders, err := payClient.GetUserOrders(c, &pb.UserOrdersRequest{
		Base:    NewBaseRequest(uc),
		UserId:  userID,
		OrderId: l.OrderID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting orders", cerr.GeneralError))
		return
	}

	if len(orders.UserOrders) != 1 {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("order_id", cerr.InvalidArgument, "No order found for id", ""))
		return
	}

	c.JSON(http.StatusOK, getRespOrder(orders.UserOrders[0]))
}

// OrderGetTransactionRequest request data
// swagger:parameters OrderGetTransactionRequest OrderGetTransaction
type OrderGetTransactionRequest struct {
	OrderID int64 `form:"order_id" json:"order_id" validate:"required"`
}

// OrderGetTransactionResponse response object
// swagger:model
type OrderGetTransactionResponse struct {
	Transaction string `json:"transaction"`
}

// OrderGetTransaction returns the unsigned payment transaction
// swagger:route GET /order_details OrderDetails
//     returns the unsigned payment transaction
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: OrderGetTransactionResponse
func OrderGetTransaction(uc *mw.IcopContext, c *gin.Context) {
	var l OrderGetTransactionRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	t, err := payClient.PayGetTransaction(c, &pb.PayGetTransactionRequest{
		Base:    NewBaseRequest(uc),
		UserId:  userID,
		OrderId: l.OrderID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting transaction", cerr.GeneralError))
		return
	}

	if t.ErrorCode != 0 {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("order_id", int(t.ErrorCode), "Error creating transaction", ""))
		return
	}

	c.JSON(http.StatusOK, &OrderGetTransactionResponse{
		Transaction: t.Transaction,
	})
}

// ExecuteTransactionRequest request data
// swagger:parameters ExecuteTransactionRequest ExecuteTransaction
type ExecuteTransactionRequest struct {
	OrderID     int64  `form:"order_id" json:"order_id" validate:"required"`
	Transaction string `form:"transaction" json:"transaction" validate:"required"`
}

// ExecuteTransactionResponse response object
// swagger:model
type ExecuteTransactionResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

// ExecuteTransaction signs the tx with the postsigner and runs the transaction
// The transaction must be signed with the customers seed on the client
// swagger:route GET /execute_transaction ExecuteTransaction
//     signs the tx with the postsigner and runs the transaction. must be signed with the customers seed on the client
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: ExecuteTransactionResponse
func ExecuteTransaction(uc *mw.IcopContext, c *gin.Context) {
	var l ExecuteTransactionRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	res, err := payClient.PayExecuteTransaction(c, &pb.PayExecuteTransactionRequest{
		Base:        NewBaseRequest(uc),
		Transaction: l.Transaction,
		OrderId:     l.OrderID,
		UserId:      userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error signing/running transaction", cerr.GeneralError))
		return
	}

	if res.ErrorCode != 0 {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("order_id", int(res.ErrorCode), "Error executing transaction", ""))
		return
	}

	c.JSON(http.StatusOK, &ExecuteTransactionResponse{
		TransactionHash: res.TxHash,
	})
}

// FakeTransactionRequest request-data
// swagger:parameters FakeTransactionRequest FakeTransaction
type FakeTransactionRequest struct {
	// required: true
	PaymentChannel string `form:"payment_channel" json:"payment_channel" validate:"required"`
	TxHash         string `form:"tx_hash" json:"tx_hash"`

	//This is the payment-address that we genrated for the order
	// required: true
	RecipientAddress string `form:"recipient_address" json:"recipient_address" validate:"required"`

	//This is the payment-address that send the payment from the external network
	// required: true
	SenderAddress string `form:"sender_address" json:"sender_address" validate:"required"`

	// required: true
	OrderID int64 `form:"order_id" json:"order_id"`

	//This is the denomination amount in the UoM of the payment network that the transaction should fake
	// required: true
	DenominationAmount int64 `form:"denomination_amount" json:"denomination_amount"`

	//This is the MEMO field for stellar payments. Must be a valid order ID
	PaymentUsage string `form:"payment_usage" json:"payment_usage"`
}

// FakeTransactionResponse response object
// swagger:model
type FakeTransactionResponse struct {
	IsDuplicate bool `json:"is_duplicate"`
}

// FakeTransaction create a fake transaction from a payment network
// swagger:route POST /fake_transaction FakeTransaction
//     create a fake transaction from a payment network
//     Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: FakeTransactionResponse
func FakeTransaction(uc *mw.IcopContext, c *gin.Context) {
	var l FakeTransactionRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if l.PaymentChannel == m.PaymentNetworkStellar && l.PaymentUsage == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("payment_usage", cerr.MissingMandatoryField, "Missing usage for stellartransaction", ""))
		return
	}

	b, err := payClient.FakePaymentTransaction(c, &pb.TestTransaction{
		Base:             NewBaseRequest(uc),
		PaymentChannel:   l.PaymentChannel,
		TxHash:           l.TxHash,
		RecipientAddress: l.RecipientAddress,
		DenomAmount:      l.DenominationAmount,
		SenderAddress:    l.SenderAddress,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating fake-transaction", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &FakeTransactionResponse{IsDuplicate: b.Value})
}
