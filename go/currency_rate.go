package java2go

type CurrencyRateMessage []CurrencyRate

type CurrencyRate struct {
	BaseCurrency   string
	QuotedCurrency string
	Rate           float64
}
