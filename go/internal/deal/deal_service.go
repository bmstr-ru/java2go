package deal

import (
	"context"
	java2go "github.com/bmstr-ru/java2go/go"
)

type DealServiceImpl struct {
	Storage         java2go.DealStorage
	ExposureService java2go.TotalExposureService
}

func (ds *DealServiceImpl) ReceiveDeal(deal *java2go.Deal) error {
	ctx := context.WithValue(context.Background(), "deal", deal)
	err := ds.Storage.SaveDeal(ctx, deal)
	if err != nil {
		return err
	}

	err = ds.ExposureService.RecalculateTotalExposure(deal.ClientId)
	if err != nil {
		return err
	}

	return nil
}
