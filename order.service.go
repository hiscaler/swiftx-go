package swiftx

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/swiftx-go/entity"
	"github.com/hiscaler/swiftx-go/response"
	"gopkg.in/guregu/null.v4"
)

// 订单服务
type orderService service

// SenderAddress 发货地址
type SenderAddress = entity.Address

// RecipientAddress 收货地址
type RecipientAddress = entity.Address

// Value 金额
type Value struct {
	Amount       float64 `json:"amount"`       // 金额数值
	CurrencyCode string  `json:"currencyCode"` // 币种代码（Enum: "USD" "CAD" "HKD" "CNY"）
}

// Validate 费用验证
func (m Value) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Amount,
			validation.Required.Error("金额不能为空"),
			validation.Min(0.0).Error("金额不能小于 0"),
		),
		validation.Field(&m.CurrencyCode,
			validation.Required.Error("币种代码不能为空"),
			is.CurrencyCode.ErrorObject(
				validation.NewError(
					"422",
					"无效的币种代码 {{.value}}").
					SetParams(map[string]interface{}{"value": m.CurrencyCode}),
			),
		),
	)
}

type CreateOrderPackageGoods struct {
	Name         string     `json:"name"`                   // SKU 名称
	Quantity     int        `json:"quantity"`               // SKU 数量
	Code         string     `json:"code,omitempty"`         // SKU 商品编码
	Value        null.Float `json:"value,omitempty"`        // SKU 单价
	CurrencyCode string     `json:"currencyCode,omitempty"` // 币种代码（Enum: "USD" "CAD" "HKD" "CNY"）
}

// Validate SKU商品验证
func (m CreateOrderPackageGoods) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name,
			validation.Required.Error("SKU 名称不能为空"),
			validation.Length(1, 255).Error("SKU 名称长度不能大于 {{.max}} 个字符"),
		),
		validation.Field(&m.Quantity,
			validation.Required.Error("SKU 数量不能为空"),
			validation.Min(1).Error("SKU 数量不能小于 {{.min}}"),
		),
		validation.Field(&m.Code,
			validation.When(m.Code != "", validation.Length(1, 100).Error("SKU 商品编码长度不能大于 {{.max}} 个字符")),
		),
		validation.Field(&m.Value, validation.When(m.Value.Valid, validation.Min(0.0).Error("SKU 单价不能小于 {{.min}}"))),
		validation.Field(&m.CurrencyCode,
			validation.When(m.CurrencyCode != "",
				is.CurrencyCode.ErrorObject(
					validation.NewError(
						"422",
						"无效的币种代码 {{.value}}").
						SetParams(map[string]any{"value": m.CurrencyCode}),
				),
			),
		),
	)
}

// CreateOrderPackageInformation 包裹信息
type CreateOrderPackageInformation struct {
	SenderAddress    SenderAddress             `json:"senderAddress"`    // 地址信息。注意：电话号码对于寄件人地址是可选的。
	RecipientAddress RecipientAddress          `json:"recipientAddress"` // 地址信息。注意：电话号码对于收件人地址是必填的。
	UseImperialUnit  bool                      `json:"useImperialUnit"`  // 是否使用英制单位，默认值false表示使用公制单位
	Weight           float64                   `json:"weight"`           // 重量(磅/公斤)
	Length           float64                   `json:"length"`           // 长度(英寸/厘米)
	Width            float64                   `json:"width"`            // 宽度(英寸/厘米)
	Height           float64                   `json:"height"`           // 高度(英寸/厘米)
	Value            Value                     `json:"value"`            // 该项费用的总额和币种
	CustomerName     string                    `json:"customerName"`     // 上游的客户名，如对接系统为ERP传入海外仓名
	StoreName        string                    `json:"storeName"`        // 电商平台店铺名
	SkuList          []CreateOrderPackageGoods `json:"skuList"`          // SKU 列表（序列化后的JSON长度不应超过 8192 字符）
}

// Validate 包裹信息验证
func (m CreateOrderPackageInformation) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SenderAddress, validation.Required.Error("发货地址不能为空"), validation.By(func(value interface{}) error {
			address, ok := value.(SenderAddress)
			if !ok {
				return errors.New("无效的发货地址")
			}
			return address.Validate()
		})),
		validation.Field(&m.RecipientAddress, validation.Required.Error("收货地址不能为空"), validation.By(func(value interface{}) error {
			address, ok := value.(RecipientAddress)
			if !ok {
				return errors.New("无效的收货地址")
			}
			if err := address.Validate(); err != nil {
				return err
			}
			if address.PhoneNumber == "" {
				return errors.New("收货地址的电话号码不能为空")
			}
			return nil
		})),
		validation.Field(&m.Weight, validation.Required.Error("重量不能为空"), validation.Min(0.0).Error("重量不能小于 0")),
		validation.Field(&m.Length, validation.Required.Error("长度不能为空"), validation.Min(0.0).Error("长度不能小于 0")),
		validation.Field(&m.Width, validation.Required.Error("宽度不能为空"), validation.Min(0.0).Error("宽度不能小于 0")),
		validation.Field(&m.Height, validation.Required.Error("高度不能为空"), validation.Min(0.0).Error("高度不能小于 0")),
		validation.Field(&m.Value, validation.Required.Error("总费用不能为空"), validation.By(func(value interface{}) error {
			v, ok := value.(Value)
			if !ok {
				return errors.New("无效的费用")
			}
			return v.Validate()
		})),
		validation.Field(&m.SkuList, validation.Required.Error("SKU 列表不能为空"), validation.Each(validation.By(func(value interface{}) error {
			return value.(CreateOrderPackageGoods).Validate()
		}))),
	)
}

// InsuranceService 保险服务配置
type InsuranceService struct {
	IsInsured    bool  `json:"isInsured"`    // 是否投保
	InsuredValue Value `json:"insuredValue"` // 该项费用的总额和币种
}

// Validate 保险服务配置验证
func (m InsuranceService) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.InsuredValue, validation.When(m.IsInsured, validation.By(func(value interface{}) error {
			v, ok := value.(Value)
			if !ok {
				return errors.New("无效的保险金额")
			}
			return v.Validate()
		}))),
	)
}

// PickupService 揽收服务配置
type PickupService struct {
	IsPickup    bool   `json:"isPickup"`    // 是否揽收
	PickupStart string `json:"pickupStart"` // 揽收开始时间
	PickupEnd   string `json:"pickupEnd"`   // 揽收结束时间
}

// Validate 揽收服务配置验证
func (m PickupService) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PickupStart, validation.When(m.IsPickup, validation.Required.Error("揽收开始时间不能为空"))),
		validation.Field(&m.PickupEnd, validation.When(m.IsPickup, validation.Required.Error("揽收结束时间不能为空"))),
	)
}

// ShippingLabelInformation 运单印刷数据
type ShippingLabelInformation struct {
	OrderNumber               string `json:"orderNumber"`                      // 上游订单号，印刷的时候会加上 "Order: " 前缀
	CustomerNote              string `json:"customerNote,omitempty"`           // 客户备注，最多 80 个字符，会印刷到快递面单上，印刷的时候会加上 "Customer note: " 前缀
	ExtSortingCode            string `json:"extSortingCode,omitempty"`         // 外部分拣码
	UseExternalTrackingNumber bool   `json:"useExternalTrackingNumber"`        // 是否使用外部面单号
	ExternalTrackingNumber    string `json:"externalTrackingNumber,omitempty"` // 外部面单号
}

// Validate 运单印刷数据验证
func (m ShippingLabelInformation) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderNumber,
			validation.Required.Error("订单号不能为空"),
			validation.Length(1, 128).Error("订单号不能超过 {{.max}} 个字符")),
		validation.Field(&m.CustomerNote, validation.When(m.CustomerNote != "", validation.Length(0, 80).Error("客户备注不能超过 {{.max}} 个字符"))),
		validation.Field(&m.ExtSortingCode, validation.When(m.ExtSortingCode != "", validation.Length(0, 128).Error("外部分拣码不能超过 {{.max}} 个字符"))),
		validation.Field(&m.ExternalTrackingNumber, validation.When(m.UseExternalTrackingNumber, validation.Length(0, 64).Error("外部面单号不能超过 {{.max}} 个字符"))),
	)
}

type CreateOrderRequest struct {
	OrderScope        string                        `json:"orderScope"`        // 订单类型,例如 DOMESTIC（国内）或 INTERNATIONAL（国际)
	ServiceType       string                        `json:"serviceType"`       // 服务类型， ECO-特惠 EXP-标快。建议选择EXP-标快
	DeliveryMethod    string                        `json:"deliveryMethod"`    // 送货方式，HDY-上门派送 SPU-自提
	CooperationMethod string                        `json:"cooperationMethod"` // 合作方式：PLATFORM-平台、MERCHANT-商家、WESTERN_POST-西邮
	ClientCode        string                        `json:"clientCode"`        // 客户代码，默认留空，使用场景需联系商务支持
	EntryPostalCode   string                        `json:"entryPostalCode"`   // 交邮点邮编，默认留空，使用场景需联系商务支持
	ReferenceNo       string                        `json:"referenceNo"`       // 引用单号，默认留空，使用场景需联系商务支持
	SelfPickupCode    string                        `json:"selfPickupCode"`    // 自提码，如果送货方式为自提时必填
	InsuranceService  *InsuranceService             `json:"insuranceService"`  // 保险服务配置
	PickupService     PickupService                 `json:"pickupService"`     // 揽收服务配置
	PackageInfo       CreateOrderPackageInformation `json:"packageInfo"`       // 包裹信息
	ShippingLabelInfo ShippingLabelInformation      `json:"shippingLabelInfo"` // 运单印刷数据
	ExtraInfo         struct {
		Platform  string `json:"platform"`
		Priority  string `json:"priority"`
		Warehouse string `json:"warehouse"`
	} `json:"extraInfo"` // 额外信息,字段可由客户自行扩展，对应字符串长度小于 4096 个字符
}

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderScope,
			validation.Required.Error("订单类型不能为空"),
			validation.In("DOMESTIC", "INTERNATIONAL").Error("无效的订单类型"),
		),
		validation.Field(&m.ServiceType,
			validation.Required.Error("服务类型不能为空"),
			validation.In("ECO", "EXP").Error("无效的服务类型"),
		),
		validation.Field(&m.DeliveryMethod,
			validation.Required.Error("送货方式不能为空"),
			validation.In("HDY", "SPU").Error("无效的送货方式"),
		),
		validation.Field(&m.CooperationMethod,
			validation.Required.Error("合作方式不能为空"),
			validation.In("PLATFORM", "MERCHANT", "WESTERN_POST").Error("无效的合作方式"),
		),
		validation.Field(&m.SelfPickupCode,
			validation.When(m.DeliveryMethod == "SPU", validation.Required.Error("自提码不能为空")),
		),
		validation.Field(&m.InsuranceService, validation.When(m.InsuranceService != nil, validation.By(func(value interface{}) error {
			v, ok := value.(*InsuranceService)
			if !ok {
				return errors.New("无效的保险服务配置")
			}
			return v.Validate()
		}))),
		validation.Field(&m.PickupService, validation.By(func(value interface{}) error {
			v, ok := value.(PickupService)
			if !ok {
				return errors.New("无效的揽收服务配置")
			}
			return v.Validate()
		})),
		validation.Field(&m.PackageInfo, validation.Required.Error("包裹信息不能为空"), validation.By(func(value interface{}) error {
			v, ok := value.(CreateOrderPackageInformation)
			if !ok {
				return errors.New("无效的包裹数据")
			}
			return v.Validate()
		})),
		validation.Field(&m.ShippingLabelInfo, validation.By(func(value interface{}) error {
			v, ok := value.(ShippingLabelInformation)
			if !ok {
				return errors.New("无效的运单印刷数据")
			}
			return v.Validate()
		})),
	)
}

type CreateOrderResult struct {
	Result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"result"`
	TrackingNo             string      `json:"trackingNo"`
	ExternalTrackingNumber null.String `json:"externalTrackingNumber,omitempty"` // 合作方跟踪号
	PdfBase64              string      `json:"pdfBase64"`
}

// Create 创建订单并获取面单 PDF 的 Base64 编码
func (s orderService) Create(ctx context.Context, request CreateOrderRequest) (entity.Order, error) {
	if err := request.Validate(); err != nil {
		return entity.Order{}, invalidInput(err)
	}

	var res CreateOrderResult
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&res).
		Post("/createOrderAndGetLabelPdfBase64")
	if err = recheckError(resp, err); err != nil {
		return entity.Order{}, err
	}
	if !res.Result.Success {
		return entity.Order{}, errors.New(res.Result.Message)
	}
	return entity.Order{
		CustomerOrderNumber: request.ShippingLabelInfo.OrderNumber,
		ShipmentNumber:      res.TrackingNo,
		TrackingNumber:      res.ExternalTrackingNumber,
		ShippingLabel:       res.PdfBase64,
	}, nil
}

// Cancel 取消订单，仅支持未揽收的订单
func (s orderService) Cancel(ctx context.Context, shipmentNumber string) (bool, error) {
	var res response.Result
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{
			"trackingNo": shipmentNumber,
		}).
		SetResult(&res).
		Post("/cancelOrder")
	if err = recheckError(resp, err); err != nil {
		return false, err
	}
	if !res.Success {
		return false, errors.New(res.Message)
	}
	return true, nil
}

// Tracking 查询物流轨迹
func (s orderService) Tracking(ctx context.Context, shipmentNumbers ...string) ([]entity.TrackingResult, error) {
	var results []entity.TrackingResult
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]string{
			"trackingNoList": shipmentNumbers,
		}).
		SetResult(&results).
		Post("/batchGetTrackingInfo")
	if err = recheckError(resp, err); err != nil {
		return nil, err
	}
	return results, nil
}
