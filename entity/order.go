package entity

import "gopkg.in/guregu/null.v4"

// Order 订单
type Order struct {
	CustomerOrderNumber string      `json:"customer_order_number"` // 客单号
	ShipmentNumber      string      `json:"shipmentNumber"`        // SwiftX 的订单号
	TrackingNumber      null.String `json:"tracking_number"`       // 上游物流商的跟踪号
	ShippingLabel       string      `json:"pdfBase64"`             // 面单 Base64
}
