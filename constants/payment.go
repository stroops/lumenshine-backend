package constants

//Chain used
type Chain string

//Currency used
type Currency string

//Constants used all over the place
const (
	/*ChainETH Chain = "ETH"
	ChainBTC Chain = "BTC"
	ChainXML Chain = "XML"

	CurrencyETH  Currency = "eth"
	CurrencyBTC  Currency = "btc"
	CurrencyXML  Currency = "xml"
	CurrencyFiat Currency = "fiat"*/

	StellarAmountPrecision = 7
)

/*
//GetChain returns the chain for the given currency
func (currency Currency) GetChain() Chain {
	if currency == CurrencyETH {
		return ChainETH
	} else if currency == CurrencyBTC {
		return ChainBTC
	} else if currency == CurrencyXML {
		return ChainXML
	}

	return ""
}

//NewChain returns a new Chain
func NewChain(chain string) Chain {
	return Chain(chain)
}

//NewCurrency returns a new Currency
func NewCurrency(currency string) Currency {
	return Currency(currency)
}

//GetCurrency returns the currency for the given chain
func (chain Chain) GetCurrency() Currency {
	if chain == ChainETH {
		return CurrencyETH
	} else if chain == ChainBTC {
		return CurrencyBTC
	} else if chain == ChainXML {
		return CurrencyXML
	}

	return CurrencyFiat //default if nothing passed
}*/
