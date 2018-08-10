package stellar

//NewAddressGenerator returns a new Addressgenerator
func NewAddressGenerator() (*AddressGenerator, error) {
	return &AddressGenerator{}, nil
}

//Generate generates a new stellar address and seed
func (g *AddressGenerator) Generate() (publicKey string, seed string, err error) {
	publicKey = "pubkey1"
	seed = "seed1"
	err = nil
	return
}
