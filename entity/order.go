package entity

// Order 订单
type Order struct {
	CustomerOrderNumber string `json:"customer_order_number"` // 客单号
	TrackingNo          string `json:"trackingNo"`            // 跟踪号
	ShippingLabel       string `json:"pdfBase64"`             // 面单 Base64
}
