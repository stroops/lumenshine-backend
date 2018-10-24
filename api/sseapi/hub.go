package main

import "fmt"

type wsMessage struct {
	//account we want to send the message to
	account string

	//message to be send
	message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	//reverse map to get the key for an address
	//dict-key is the address, value is a list of all client-keys this address is registered for
	addresses map[string][]string

	send chan *wsMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type address string
type key string

func newHub() *Hub {
	return &Hub{
		send: make(chan *wsMessage),

		register:   make(chan *Client),
		unregister: make(chan *Client),
		//clients:    make(map[*Client]bool),

		clients: make(map[string]*Client),

		//reverse map to get the key for an address
		//dict-key is the address, value the key for reverse-lookup of the client
		addresses: make(map[string][]string),
	}
}

func (h *Hub) removeClient(client *Client) {

	//remove the client from the hub address - list
	addresses, _ := h.addresses[client.key]

	for i := len(addresses) - 1; i >= 0; i-- {
		for _, ca := range client.addresses {
			if addresses[i] == ca {
				addresses = append(addresses[:i], addresses[i+1:]...)
			}
		}
	}
	if len(addresses) == 0 {
		delete(h.addresses, client.key)
	}

	//close and delete the client
	client.conn.Close()
	close(client.send)
	delete(h.clients, client.key)
	fmt.Println("Closed client")
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
			keys, ok := h.addresses[message.account]
			if ok {
				//send message to all clients for this account
				for _, key := range keys {
					client, ok := h.clients[key]
					if ok {
						select {
						case client.send <- message.message:
						default:
							h.removeClient(client)
						}
					}
				}
			}
		}
	}
}
