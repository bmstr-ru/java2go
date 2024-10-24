package java2go

type TotalExposure struct {
	ClientId int64            `json:"clientId"`
	Total    MonetaryAmount   `json:"total"`
	Amounts  []MonetaryAmount `json:"amounts"`
}

type TotalExposureService interface {
	RecalculateAllTotalExposure() error
	GetClientsTotalExposure(clientId int64) (*TotalExposure, error)
	ConsiderNewAmounts(clientId int64, monetaryAmounts ...MonetaryAmount) error
}
