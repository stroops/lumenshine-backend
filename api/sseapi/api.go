package main

import (
	"net/http"
	"time"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"
	"github.com/Soneso/lumenshine-backend/helpers"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"
	"github.com/gin-gonic/gin"
)

var bits helpers.Bits

func init() {
	bits = helpers.Set(bits, helpers.F0) //create
	bits = helpers.Set(bits, helpers.F1) //payment
	bits = helpers.Set(bits, helpers.F2) //paymentPath
}

//GetWSRequest - requestdata for a websocket
//swagger:parameters GetWSRequest GetWS
type GetWSRequest struct {
	//Random key for identlifying the websocket on the backend.
	//required: true
	RandomKey string `form:"random_key" json:"random_key" query:"random_key" validate:"required"`
}

// GetWS returns a new websocket connection
// swagger:route GET /portal/sse/get_ws websocket GetWS
//
// Returns a new websocket connection
// 	  Consumes:
//     - multipart/form-data
//
//     Responses:
//       200:
func GetWS(hub *Hub, uc *mw.IcopContext, c *gin.Context) {
	var l GetWSRequest
	if !helpers.ValidateRequestData(&l, uc.Log, c) {
		return
	}

	//check that key not used yet
	if _, ok := hub.clients[l.RandomKey]; ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("random_key", cerr.InvalidArgument, "Key already exists", ""))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("random_key", cerr.InvalidArgument, "Upgrader not created", ""))
		return
	}

	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 1024), key: l.RandomKey}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
}

//RemoveWSRequest - requestdata for removing a websocket
//swagger:parameters RemoveWSRequest RemoveWS
type RemoveWSRequest struct {
	// RandomKey from the client, that identifies the websocket. Was provided by get_ws
	// required: true
	Key string `form:"key" json:"key" query:"key" validate:"required"`
}

// RemoveWS removes a Websocket connection
// swagger:route POST /portal/sse/remove_ws websocket RemoveWS
//
// Removes a Websocket connection
// 	  Consumes:
//     - multipart/form-data
//
//     Responses:
//       200:
func RemoveWS(hub *Hub, uc *mw.IcopContext, c *gin.Context) {
	var l RemoveWSRequest
	if !helpers.ValidateRequestData(&l, uc.Log, c) {
		return
	}

	client, ok := hub.clients[l.Key]
	if !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("key", cerr.InvalidArgument, "Key not found", ""))
		return
	}

	hub.removeClient(client)
	c.JSON(http.StatusOK, "{}")
}

//ListenAccountRequest requestdata for adding an account listener
//swagger:parameters ListenAccountRequest ListenAccount
type ListenAccountRequest struct {
	// RandomKey from the client, that identifies the websocket. Was provided by get_ws
	// required: true
	Key string `form:"key" json:"key" query:"key" validate:"required"`

	// Account is the stellar account the client wants to listen for events
	// required: true
	Account string `form:"account" json:"account" query:"account" validate:"required"`
}

// ListenAccount adds a listener for an account
// swagger:route POST /portal/sse/listen_account websocket ListenAccount
//
// Adds a listener for an account
// 	  Consumes:
//     - multipart/form-data
//
//     Responses:
//       200:
func ListenAccount(hub *Hub, uc *mw.IcopContext, c *gin.Context) {
	var l ListenAccountRequest
	if !helpers.ValidateRequestData(&l, uc.Log, c) {
		return
	}

	client, ok := hub.clients[l.Key]
	if !ok {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("key", cerr.InvalidArgument, "Key not found", ""))
		return
	}

	//check if address already  registered for key
	for _, a := range client.addresses {
		if a == l.Account {
			c.JSON(http.StatusOK, "{}")
			return
		}
	}

	//register account in sse for events
	_, err := sseClient.ListenFor(c, &pb.SSEListenForRequest{
		Base:           NewBaseRequest(uc),
		OpTypes:        int64(bits),
		SourceReciver:  "sse",
		StellarAccount: l.Account,
		WithResume:     false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error registering account for sse", cerr.GeneralError))
		return
	}

	//not found add address to listener and reverselookup
	client.addresses = append(client.addresses, l.Account)
	hub.addresses[l.Account] = append(hub.addresses[l.Account], l.Key)
	c.JSON(http.StatusOK, "{}")
}

//RemoveAccountRequest requestdata for removing an account listener
//swagger:parameters RemoveAccountRequest RemoveAccount
type RemoveAccountRequest struct {
	// RandomKey from the client, that identifies the websocket. Was provided by get_ws
	// required: true
	Key string `form:"key" json:"key" query:"key" validate:"required"`

	// Account is the stellar account the client wants to remove the event listening
	// required: true
	Account string `form:"account" json:"account" query:"account" validate:"required"`
}

// RemoveAccount removes a listener for an account
// swagger:route POST /portal/sse/remove_account websocket RemoveAccount
//
// Removes a listener for an account
// 	  Consumes:
//     - multipart/form-data
//
//     Responses:
//       200:
func RemoveAccount(hub *Hub, uc *mw.IcopContext, c *gin.Context) {
	var l RemoveAccountRequest
	if !helpers.ValidateRequestData(&l, uc.Log, c) {
		return
	}

	if client, ok := hub.clients[l.Key]; ok {
		hub.removeAddress(client, l.Account)
	}

	c.JSON(http.StatusOK, "{}")
}

//SendMessageRequest requestdata for sending data to an client
//swagger:parameters SendMessageRequest SendMessage
type SendMessageRequest struct {
	// required: true
	Account string `form:"account" json:"account" query:"account" validate:"required"`

	// required: true
	MessageType int64 `form:"message_type" json:"message_type" query:"message_type" validate:"required"`

	// required: true
	Data string `form:"data" json:"data" query:"data"`
}

// SendMessage send a message to all clients listening for the account
// swagger:route POST /send_message websocket SendMessage
//
// Send a message to all clients listening for the account. Can only be used when in development mode
// 	  Consumes:
//     - multipart/form-data
//
//     Responses:
//       200:
func SendMessage(hub *Hub, uc *mw.IcopContext, c *gin.Context) {
	var l SendMessageRequest
	if !helpers.ValidateRequestData(&l, uc.Log, c) {
		return
	}

	_, ok := hub.addresses[l.Account]
	if ok {
		hub.send <- &WsMessage{
			Account:     l.Account,
			MessageType: l.MessageType,
			Message:     []byte(l.Data),
		}
	}

	c.JSON(http.StatusOK, "{}")
}
