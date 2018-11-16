package main

import (
	"encoding/json"
	"fmt"
	"sync"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/pb"
	context "golang.org/x/net/context"
)

//WsMessage is a message to be send to the websocket
type WsMessage struct {
	//account we want to send the message to
	Account string `json:"account"`

	//MessageType int64 `json:"message_type"`

	//message to be send
	//Message []byte `json:"message"`
}

//GetJSON returns the json string for a message
func (w *WsMessage) GetJSON() []byte {
	b, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
		return []byte("{}")
	}
	return b
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients. MapKey is the ranom client key
	clients map[string]*Client

	//reverse map to get the key for an address
	//dict-key is the address, value is a list of all client-keys this address is registered for
	addresses map[string][]string

	send chan *WsMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mux sync.Mutex
}

type address string
type key string

func newHub() *Hub {
	//  Clear current data for sse
	ctx := context.Background()
	_, err := sseClient.ClearSourceRecivers(ctx, &pb.SSEClearSourceReciversRequest{
		Base:          &pb.BaseRequest{RequestId: "0", UpdateBy: ServiceName},
		SourceReciver: m.SourceReceiverSse,
	})
	if err != nil {
		fmt.Printf("Error deleting sse-data %v", err)
	}

	return &Hub{
		send: make(chan *WsMessage),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		//clients:    make(map[*Client]bool),

		clients: make(map[string]*Client),

		//reverse map to get the keys for an address
		//dict-key is the address, value the key for reverse-lookup of the client
		addresses: make(map[string][]string),
	}
}

func (h *Hub) removeAddress(client *Client, account string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	//delete address from client-list
	for i, a := range client.addresses {
		if a == account {
			copy(client.addresses[i:], client.addresses[i+1:])
			client.addresses[len(client.addresses)-1] = ""
			client.addresses = client.addresses[:len(client.addresses)-1]
		}
	}

	//remove key from addresses lookup
	keys, ok := h.addresses[account]
	if ok {
		for i := 0; i < len(keys); i++ {
			if keys[i] == client.key {
				copy(keys[i:], keys[i+1:])
				keys[len(keys)-1] = ""
				keys = keys[:len(keys)-1]
				h.addresses[account] = keys
			}
		}
	}

	if len(keys) == 0 {
		delete(h.addresses, account)
		fmt.Printf("Deleting account %s", account)

		//if we do not listen for the address any longer, we also need to remove the address from the sse-service
		//register account in sse for events
		ctx := context.Background()
		_, err := sseClient.RemoveListening(ctx, &pb.SSERemoveListeningRequest{
			Base:           &pb.BaseRequest{RequestId: "0", UpdateBy: ServiceName},
			SourceReciver:  m.SourceReceiverSse,
			StellarAccount: account,
		})
		if err != nil {
			fmt.Printf("Error deleting sse-account %v", err)
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[client.key]; ok {
		for _, ca := range client.addresses {
			hub.removeAddress(client, ca)
		}

		//close and delete the client
		client.conn.Close()
		close(client.send)
		delete(h.clients, client.key)
		fmt.Println("Closed client")
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			_, ok := h.clients[client.key]
			if !ok {
				//does not exist yet, add client
				h.clients[client.key] = client
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.key]; ok {
				fmt.Println("Closing unregister")
				h.removeClient(client)
			}
		case message := <-h.send:
			keys, ok := h.addresses[message.Account]
			if ok {
				//send message to all clients for this account
				for _, key := range keys {
					client, ok := h.clients[key]
					if ok {
						select {
						case client.send <- message.GetJSON():
						default:
							h.removeClient(client)
						}
					}
				}
			}
		}
	}
}
