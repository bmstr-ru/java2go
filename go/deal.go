package java2go

import "fmt"

type Deal struct {
	Id           int64          `json:"id"`
	ClientId     int64          `json:"clientId"`
	AmountBought MonetaryAmount `json:"amountBought"`
	AmountSold   MonetaryAmount `json:"amountSold"`
}

func (d *Deal) String() string {
	return fmt.Sprintf("deal[id=%d, clientId=%d, amountBought=%s, amountSold=%s]",
		d.Id, d.ClientId, d.AmountBought.String(), d.AmountSold.String())
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

func (m *MonetaryAmount) String() string {
	return fmt.Sprintf("%f %s", m.Amount, m.Currency)
}

type DealService interface {
	ReceiveDeal(deal *Deal) error
}
