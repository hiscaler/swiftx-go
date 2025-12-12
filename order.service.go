package swiftx

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/swiftx-go/entity"
)

// 订单服务
type orderService service

// SenderAddress 发货地址
type SenderAddress = entity.Address

// RecipientAddress 收货地址
type RecipientAddress = entity.Address

type CreateOrderPackageGoods struct {
	Name         string  `json:"name"`         // SKU 名称
	Quantity     int     `json:"quantity"`     // SKU 数量
	Code         string  `json:"code"`         // SKU 商品编码
	Value        float64 `json:"value"`        // SKU 单价
	CurrencyCode string  `json:"currencyCode"` // 币种代码（Enum: "USD" "CAD" "HKD" "CNY"）
}

// Value 金额
type Value struct {
	Amount       float64 `json:"amount"`       // 金额数值
	CurrencyCode string  `json:"currencyCode"` // 币种代码（Enum: "USD" "CAD" "HKD" "CNY"）
}

// CreateOrderPackageInformation 包裹信息
type CreateOrderPackageInformation struct {
	SenderAddress    SenderAddress             `json:"senderAddress"`    // 地址信息。注意：电话号码对于收件人地址是必填的，对于寄件人地址是可选的。
	RecipientAddress RecipientAddress          `json:"recipientAddress"` // 地址信息。注意：电话号码对于收件人地址是必填的，对于寄件人地址是可选的。
	UseImperialUnit  bool                      `json:"useImperialUnit"`  // 是否使用英制单位，默认值false表示使用公制单位
	Weight           float64                   `json:"weight"`           // 重量(磅/公斤)
	Length           int                       `json:"length"`           // 长度(英寸/厘米)
	Width            int                       `json:"width"`            // 宽度(英寸/厘米)
	Height           float64                   `json:"height"`           // 高度(英寸/厘米)
	Value            Value                     `json:"value"`            // 该项费用的总额和币种
	CustomerName     string                    `json:"customerName"`     // 上游的客户名，如对接系统为ERP传入海外仓名
	StoreName        string                    `json:"storeName"`        // 电商平台店铺名
	SkuList          []CreateOrderPackageGoods `json:"skuList"`          // SKU列表（序列化后的JSON长度不应超过8192字符）
}

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
		IsInsured    bool  `json:"isInsured"`    // 是否投保
		InsuredValue Value `json:"insuredValue"` // 该项费用的总额和币种
	} `json:"insuranceService"` // 保险服务配置
	PickupService struct {
		IsPickup    bool   `json:"isPickup"`    // 是否揽收
		PickupStart string `json:"pickupStart"` // 揽收开始时间
		PickupEnd   string `json:"pickupEnd"`   // 揽收结束时间
	} `json:"pickupService"`                                             // 揽收服务配置
	PackageInfo       CreateOrderPackageInformation `json:"packageInfo"` // 包裹信息
	ShippingLabelInfo struct {
		OrderNumber               string `json:"orderNumber"`               // 上游订单号，印刷的时候会加上 "Order: " 前缀
		CustomerNote              string `json:"customerNote"`              // 客户备注，最多80个字符，会印刷到快递面单上，印刷的时候会加上 "Customer note: " 前缀
		ExtSortingCode            string `json:"extSortingCode"`            // 外部分拣码
		UseExternalTrackingNumber bool   `json:"useExternalTrackingNumber"` // 是否使用外部面单号
		ExternalTrackingNumber    string `json:"externalTrackingNumber"`    // 外部面单号
	} `json:"shippingLabelInfo"` // 运单印刷数据
	ExtraInfo struct {
		Platform  string `json:"platform"`
		Priority  string `json:"priority"`
		Warehouse string `json:"warehouse"`
	} `json:"extraInfo"` // 额外信息,字段可由客户自行扩展，对应字符串长度小于4096个字符
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
