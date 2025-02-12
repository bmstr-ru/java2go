package java2go

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestPrintDeal(t *testing.T) {
	deal := Deal{
		Id:           1,
		ClientId:     2,
		AmountBought: MonetaryAmount{Amount: 32.123, Currency: "USD"},
		AmountSold:   MonetaryAmount{Amount: 482.543, Currency: "EUR"},
	}
	log.Info().Str("deal", deal.String()).Msg("Test String()")
}
