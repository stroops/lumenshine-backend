package main

import (
	"fmt"
)

//GeneratePaymentAddress generates a new address in the payment network
//checks if the generator allredy exists and if not, creates one
func (s *server) GeneratePaymentAddress(paymentNetwork string, masterKey string) (address string, privatekey string, err error) {
	client, ok := s.Env.Clients[paymentNetwork]
	if !ok {
		err = fmt.Errorf("payment-network %s does not exist", paymentNetwork)
		return
	}

	address, privatekey, err = client.GeneratePaymentAddress()
	return
}
