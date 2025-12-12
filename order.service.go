package swiftx

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/swiftx-go/entity"
)

// 订单服务
type orderService service

type CreateOrderRequest struct {
	OrderScope        string `json:"orderScope"`        // 订单类型,例如 DOMESTIC（国内）或 INTERNATIONAL（国际)
	ServiceType       string `json:"serviceType"`       // 服务类型， ECO-特惠 EXP-标快。建议选择EXP-标快
	DeliveryMethod    string `json:"deliveryMethod"`    // 送货方式，HDY-上门派送 SPU-自提
	CooperationMethod string `json:"cooperationMethod"` // 合作方式：PLATFORM-平台、MERCHANT-商家、WESTERN_POST-西邮
	ClientCode        string `json:"clientCode"`        // 客户代码，默认留空，使用场景需联系商务支持
	EntryPostalCode   string `json:"entryPostalCode"`   // 交邮点邮编，默认留空，使用场景需联系商务支持
	ReferenceNo       string `json:"referenceNo"`       // 引用单号，默认留空，使用场景需联系商务支持
	SelfPickupCode    string `json:"selfPickupCode"`    // 自提码，如果送货方式为自提时必填
	InsuranceService  struct {
		IsInsured    bool `json:"isInsured"` // 是否投保
		InsuredValue struct {
			Amount       int    `json:"amount"`       // 金额数值
			CurrencyCode string `json:"currencyCode"` // 币种代码（Enum: "USD" "CAD" "HKD" "CNY"）
		} `json:"insuredValue"` // 该项费用的总额和币种
	} `json:"insuranceService"` // 保险服务配置
	PickupService struct {
		IsPickup    bool   `json:"isPickup"`    // 是否揽收
		PickupStart string `json:"pickupStart"` // 揽收开始时间
		PickupEnd   string `json:"pickupEnd"`   // 揽收结束时间
	} `json:"pickupService"` // 揽收服务配置
	PackageInfo struct {
		SenderAddress struct {
			Name             string `json:"name"`             // 姓名
			PhoneCountryCode int    `json:"phoneCountryCode"` // 电话国际区号。如果留空，将使用regionCode对应的国际区号作为电话的区号
			PhoneNumber      string `json:"phoneNumber"`      // 电话号码（寄件人地址可选，收件人地址必填）
			Building         string `json:"building"`         // 电话分机号
			StreetAddress    string `json:"streetAddress"`    // 建筑物名称，对应美国地址的address line 2
			District         string `json:"district"`         // 街道地址，对应美国地址的address line 1
			City             string `json:"city"`             // 城市
			StateProvince    string `json:"stateProvince"`    // 州/省（美国州请使用2字母缩写，如CA、NY）
			PostalCode       string `json:"postalCode"`       // 邮政编码，对应美国的zipcode。最大支持10个字符。美国zipcode支持标准5位格式（如94105）或带+4扩展代码的格式（如94105-1234）。当使用+4格式时，前5位用于业务校验和路线规划，后4位仅作为记录保存在订单数据中。SwiftX服务覆盖邮编，请联系商务获取。
			RegionCode       string `json:"regionCode"`       // 国家编码，类似US、CN（Enum: "US" "CA" "CN"）
		} `json:"senderAddress"` // 地址信息。注意：电话号码对于收件人地址是必填的，对于寄件人地址是可选的。
		RecipientAddress struct {
			Name             string `json:"name"`             // aaaa
			PhoneCountryCode int    `json:"phoneCountryCode"` // aaaa
			PhoneNumber      string `json:"phoneNumber"`      // aaaa
			PhoneExtension   string `json:"phoneExtension"`   // aaaa
			Building         string `json:"building"`         // aaaa
			StreetAddress    string `json:"streetAddress"`    // aaaa
			District         string `json:"district"`         // aaaa
			City             string `json:"city"`             // aaaa
			StateProvince    string `json:"stateProvince"`    // aaaa
			PostalCode       string `json:"postalCode"`       // aaaa
			RegionCode       string `json:"regionCode"`       // aaaa
		} `json:"recipientAddress"`                      // aaaa
		UseImperialUnit bool    `json:"useImperialUnit"` // aaaa
		Weight          float64 `json:"weight"`          // aaaa
		Length          int     `json:"length"`          // aaaa
		Width           int     `json:"width"`           // aaaa
		Height          float64 `json:"height"`          // aaaa
		Value           struct {
			Amount       float64 `json:"amount"`
			CurrencyCode string  `json:"currencyCode"`
		} `json:"value"`
		CustomerName string `json:"customerName"`
		StoreName    string `json:"storeName"`
		SkuList      []struct {
			Name         string  `json:"name"`
			Quantity     int     `json:"quantity"`
			Code         string  `json:"code"`
			Value        float64 `json:"value"`
			CurrencyCode string  `json:"currencyCode"`
		} `json:"skuList"`
	} `json:"packageInfo"` // 包裹信息
	ShippingLabelInfo struct {
		OrderNumber    string `json:"orderNumber"`
		CustomerNote   string `json:"customerNote"`
		ExtSortingCode string `json:"extSortingCode"`
	} `json:"shippingLabelInfo"`
	ExtraInfo struct {
		Platform  string `json:"platform"`
		Priority  string `json:"priority"`
		Warehouse string `json:"warehouse"`
	} `json:"extraInfo"`
}

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderScope,
			validation.Required.Error("订单类型不能为空"),
			validation.In("DOMESTIC", "INTERNATIONAL").ErrorObject(validation.NewError("422", "无效的订单类型 {{.value}}").SetParams(map[string]interface{}{"value": m.OrderScope})),
		),
	)
}

type CreateOrderResult struct {
	Result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"result"`
	TrackingNo string `json:"trackingNo"`
	PdfBase64  string `json:"pdfBase64"`
}

// Create 创建订单并获取面单 PDF 的 Base64 编码
func (s orderService) Create(ctx context.Context, requests []CreateOrderRequest) ([]CreateOrderResult, error) {
	for _, req := range requests {
		if err := req.Validate(); err != nil {
			return nil, invalidInput(err)
		}
	}

	var res []CreateOrderResult
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(requests).
		SetResult(&res).
		Post("/createOrderAndGetLabelPdfBase64")
	if err = recheckError(resp, err); err != nil {
		return nil, err
	}
	return res, nil
}

// Cancel 取消订单，仅支持未揽收的订单
func (s orderService) Cancel(ctx context.Context, trackingNumber string) (bool, error) {
	var res struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{
			"trackingNo": trackingNumber,
		}).
		SetResult(&res).
		Post("/cancelOrder")
	if err = recheckError(resp, err); err != nil {
		return false, err
	}
	return res.Success, nil
}

// Tracking 查询物流轨迹
func (s orderService) Tracking(ctx context.Context, trackingNumbers ...string) ([]entity.TrackingResult, error) {
	var res []entity.TrackingResult
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]string{
			"trackingNoList": trackingNumbers,
		}).
		SetResult(&res).
		Post("/batchGetTrackingInfo")
	if err = recheckError(resp, err); err != nil {
		return nil, err
	}
	return res, nil
}
