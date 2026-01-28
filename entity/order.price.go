package entity

// 订单价格

type Money struct {
	CurrencyCode string  `json:"currencyCode"`
	Value        float64 `json:"value"`
}

type PriceDetail struct {
	Cost        Money  `json:"cost"`
	Description string `json:"description"`
}
type OrderPrice struct {
	TrackingNumber string        `json:"tracking_number"`
	Amount         Money         `json:"amount"`
	Details        []PriceDetail `json:"details"`
}
