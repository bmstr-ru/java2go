package java2go

type TotalExposure struct {
	ClientId int64
	Total    MonetaryAmount
	Amounts  []MonetaryAmount
}

type TotalExposureService interface {
	RecalculateAllTotalExposure() error
	RecalculateTotalExposure(clientId int64) error
	GetClientsTotalExposure(clientId int64) *TotalExposure
	ConsiderNewAmounts(clientId int64, monetaryAmounts []MonetaryAmount) error
}
