package java2go

type CurrencyRateMessage []CurrencyRate

type CurrencyRate struct {
	BaseCurrency   string  `json:"baseCurrency"`
	QuotedCurrency string  `json:"quotedCurrency"`
	Rate           float64 `json:"rate"`
}

type CurrencyRateService interface {
	ReceiveRate(rate *CurrencyRate) error
}
