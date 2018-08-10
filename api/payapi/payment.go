package main

import (
	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/pb"
	"net/http"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	m "github.com/Soneso/lumenshine-backend/services/db/models"

	"github.com/gin-gonic/gin"
)

//PriceForCoinRequest struct
type PriceForCoinRequest struct {
	CoinAmount int64  `form:"coin_amount" json:"coin_amount" validate:"required"`
	Chain      string `form:"chain" json:"chain" validate:"required,oneof=eth btc xlm fiat"`
}

//PriceForCoinResponse return struct
type PriceForCoinResponse struct {
	ChainAmount             float64 `json:"chain_amount"`
	ChainAmountDenomination string  `json:"chain_amount_denomination"`
}

//PriceForCoin returns the price for a given coin amount, related to the currency given
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

	// TODO: Checks for payment
	// correct state,amount, etc.
	price, err := payClient.GetCoinPrice(c, &pb.CoinPriceRequest{
		Base:       NewBaseRequest(uc),
		Chain:      l.Chain,
		CoinAmount: l.CoinAmount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting price", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &PriceForCoinResponse{
		ChainAmount:             price.ChainAmount,
		ChainAmountDenomination: price.ChainAmountDenomination,
	})
}

//CreateOrderRequest is the data needed for tmaking a payment
type CreateOrderRequest struct {
	CoinAmount    int64  `form:"coin_amount" json:"coin_amount" validate:"required"`
	Chain         string `form:"chain" json:"chain" validate:"required,oneof=xlm eth btc fiat"`
	UserPublicKey string `form:"user_public_key"`
}

//CreateOrderResponse return struct
type CreateOrderResponse struct {
	OrderID             int64   `json:"order_id"`
	ChainAddress        string  `json:"chain_address,omitempty"`
	Chain               string  `json:"chain"`
	ChainAmount         float64 `json:"chain_amount"`
	FiatBIC             string  `json:"fiat_bic,omitempty"`
	FiatIBAN            string  `json:"fiat_iban,omitempty"`
	FiatDestinationName string  `json:"fiat_destination_name,omitempty"`
	FiatPaymentUsage    string  `json:"fiat_payment_usage,omitempty"`
}

//CreateOrder called from the API
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

	phase, err := payClient.GetActveICOPhase(c, &pb.Empty{Base: NewBaseRequest(uc)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting ico phase", cerr.GeneralError))
		return
	}

	if phase == nil || phase.PhaseName == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("Chain", cerr.NoActivePhase, "No active Phase", ""))
		return
	}

	if phase.CoinAmount < l.CoinAmount {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("Chain", cerr.InsufficientCoins, "Insuffitiant coins for phase", ""))
		return
	}

	// TODO: Checks for payment
	// correct amount, ordercount, check for fiat, etc.
	if l.Chain != m.ChainFiat && l.UserPublicKey == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("user_public_key", cerr.MissingMandatoryField, "Need stellar public key", ""))
		return
	}

	orderReq := &pb.CreateOrderRequest{
		Base:          NewBaseRequest(uc),
		UserId:        mw.GetAuthUser(c).UserID,
		Chain:         l.Chain,
		UserPublicKey: l.UserPublicKey,
		CoinAmount:    l.CoinAmount,
		IcoPhase:      phase.PhaseName,
	}
	order, err := payClient.CreateOrder(c, orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error creating order", cerr.GeneralError))
		return
	}

	var r CreateOrderResponse
	r.Chain = order.Chain
	r.ChainAmount = order.ChainAmount
	r.ChainAddress = order.Address
	r.OrderID = order.OrderId
	r.FiatBIC = order.FiatBic
	r.FiatDestinationName = order.FiatDestinationName
	r.FiatIBAN = order.FiatIban
	r.FiatPaymentUsage = order.FiatPaymentUsage

	c.JSON(http.StatusOK, r)
}

//UserOrderResponse resp object
type UserOrderResponse struct {
	ID                   int64   `json:"id"`
	OrderStatus          string  `json:"order_status"`
	CoinAmount           int64   `json:"coin_amount"`
	ChainAmount          float64 `json:"chain_amount"`
	ChainAmountDenom     string  `json:"chain_amount_denom"`
	Chain                string  `json:"chain,omitempty"`
	ChainAddress         string  `json:"chain_address,omitempty"`
	UserStellarPublicKey string  `json:"user_stellar_public_key,omitempty"`
	FiatBic              string  `json:"fiat_bic,omitempty"`
	FiatIban             string  `json:"fiat_iban,omitempty"`
	FiatDestinationName  string  `json:"fiat_destination_name,omitempty"`
	FiatPaymentUsage     string  `json:"fiat_payment_usage,omitempty"`
}

//OrderList returns the list of all orders for the user
func OrderList(uc *mw.IcopContext, c *gin.Context) {
	userID := mw.GetAuthUser(c).UserID
	orders, err := payClient.GetUserOrders(c, &pb.UserOrdersRequest{
		Base:   NewBaseRequest(uc),
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting orders", cerr.GeneralError))
		return
	}

	ret := make([]UserOrderResponse, len(orders.UserOrders))
	for i := 0; i < len(orders.UserOrders); i++ {
		o := orders.UserOrders[i]
		ret[i].ID = o.Id
		ret[i].OrderStatus = o.OrderStatus
		ret[i].CoinAmount = o.CoinAmount
		ret[i].ChainAmount = o.ChainAmount
		ret[i].ChainAmountDenom = o.ChainAmountDenom
		ret[i].Chain = o.Chain
		ret[i].ChainAddress = o.ChainAddress
		ret[i].UserStellarPublicKey = o.UserStellarPublicKey
		ret[i].FiatBic = o.FiatBic
		ret[i].FiatIban = o.FiatIban
		ret[i].FiatDestinationName = o.FiatDestinationName
		ret[i].FiatPaymentUsage = o.FiatPaymentUsage
	}

	c.JSON(http.StatusOK, ret)
}

//OrderDetailsRequest request-data
type OrderDetailsRequest struct {
	OrderID int64 `form:"order_id" json:"order_id" validate:"required"`
}

//OrderDetails returns the detail of one order
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

	o := orders.UserOrders[0]
	c.JSON(http.StatusOK, &UserOrderResponse{
		ID:                   o.Id,
		OrderStatus:          o.OrderStatus,
		CoinAmount:           o.CoinAmount,
		ChainAmount:          o.ChainAmount,
		ChainAmountDenom:     o.ChainAmountDenom,
		Chain:                o.Chain,
		ChainAddress:         o.ChainAddress,
		UserStellarPublicKey: o.UserStellarPublicKey,
		FiatBic:              o.FiatBic,
		FiatIban:             o.FiatIban,
		FiatDestinationName:  o.FiatDestinationName,
		FiatPaymentUsage:     o.FiatPaymentUsage,
	})
}

//SetOrderStatusRequest request-data
type SetOrderStatusRequest struct {
	OrderID      int64  `form:"order_id" json:"order_id" validate:"required"`
	Status       string `form:"status" json:"status" validate:"required,oneof=payment_error finished"`
	ErrorMessage string `form:"error_message" json:"error_message"`
}

//SetOrderStatus sets the status of the order
//can only be used, if the order is in status tx_created
func SetOrderStatus(uc *mw.IcopContext, c *gin.Context) {
	var l SetOrderStatusRequest
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

	order := orders.UserOrders[0]
	if order.OrderStatus != m.OrderStatusTXCreated {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("order_id", cerr.OrderWrongStatus, "Order has wrong status", ""))
		return
	}

	_, err = payClient.UserOrderSetStatus(c, &pb.UserOrderSetStatusRequest{
		Base:         NewBaseRequest(uc),
		UserId:       userID,
		OrderId:      l.OrderID,
		Status:       l.Status,
		ErrorMessage: l.ErrorMessage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error setting order-status", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}

//OrderGetTrustStatusRequest request data
type OrderGetTrustStatusRequest struct {
	OrderID int64 `form:"order_id" json:"order_id" validate:"required"`
}

//OrderGetTrustStatusResponse response onbect
type OrderGetTrustStatusResponse struct {
	HasTrustline bool `json:"has_trustline"`
}

//OrderGetTrustStatus returns the status for the trustline
//also creates the stellar user account, if it does not exist
func OrderGetTrustStatus(uc *mw.IcopContext, c *gin.Context) {
	var l OrderGetTrustStatusRequest
	if err := c.Bind(&l); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, l); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	userID := mw.GetAuthUser(c).UserID

	t, err := payClient.PayGetTrustStatus(c, &pb.PayGetTrustStatusRequest{
		Base:    NewBaseRequest(uc),
		UserId:  userID,
		OrderId: l.OrderID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error getting trust status", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &OrderGetTrustStatusResponse{
		HasTrustline: t.Value,
	})
}

//OrderGetTransactionRequest request data
type OrderGetTransactionRequest struct {
	OrderID int64 `form:"order_id" json:"order_id" validate:"required"`
}

//OrderGetTransactionResponse response onbect
type OrderGetTransactionResponse struct {
	Transaction string `json:"transaction"`
}

//OrderGetTransaction returns the payment transaction unsigned
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

//ExecuteTransactionRequest -
type ExecuteTransactionRequest struct {
	OrderID     int64  `form:"order_id" json:"order_id" validate:"required"`
	Transaction string `form:"transaction" json:"transaction" validate:"required"`
}

//ExecuteTransaction signs and runs a transaction
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

	_, err := payClient.PayExecuteTransaction(c, &pb.PayExecuteTransactionRequest{
		Base:        NewBaseRequest(uc),
		Transaction: l.Transaction,
		OrderId:     l.OrderID,
		UserId:      userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error signing/running transaction", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, "{}")
}
