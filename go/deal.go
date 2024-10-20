package java2go

type Deal struct {
	Id           int64          `json:"id"`
	ClientId     int64          `json:"clientId"`
	AmountBought MonetaryAmount `json:"amountBought"`
	AmountSold   MonetaryAmount `json:"amountSold"`
}

type MonetaryAmount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

func (m *MonetaryAmount) Negate() MonetaryAmount {
	return MonetaryAmount{
		Currency: m.Currency,
		Amount:   -m.Amount,
	}
}

type DealService interface {
	ReceiveDeal(deal *Deal) error
}
