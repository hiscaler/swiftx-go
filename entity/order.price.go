package entity

// 订单价格

type Money struct {
	CurrencyCode string  `json:"currencyCode"` // 币种代码
	Value        float64 `json:"value"`        // 金额
}

type PriceDetail struct {
	Cost        Money  `json:"cost"`        // 金额
	Description string `json:"description"` // 描述
}
type OrderPrice struct {
	TrackingNumber string        `json:"tracking_number"` // 跟踪号
	Amount         Money         `json:"amount"`          // 金额
	Details        []PriceDetail `json:"details"`         // 详情
}
