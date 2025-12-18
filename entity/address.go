package entity

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Address 收发货地址
type Address struct {
	Name             string `json:"name"`                       // 姓名
	PhoneCountryCode int    `json:"phoneCountryCode,omitempty"` // 电话国际区号。如果留空，将使用regionCode对应的国际区号作为电话的区号
	PhoneNumber      string `json:"phoneNumber,omitempty"`      // 电话号码（寄件人地址可选，收件人地址必填）
	PhoneExtension   string `json:"phoneExtension,omitempty"`   // 电话分机号
	RegionCode       string `json:"regionCode"`                 // 国家编码，类似US、CN（Enum: "US" "CA" "CN"）
	StateProvince    string `json:"stateProvince"`              // 省/州（美国州请使用2字母缩写，如CA、NY）
	City             string `json:"city"`                       // 城市
	District         string `json:"district,omitempty"`         // 区域/县，对应美国地址的county
	StreetAddress    string `json:"streetAddress"`              // 街道地址，对应美国地址的address line 1
	Building         string `json:"building,omitempty"`         // 建筑物名称，对应美国地址的address line 2
	PostalCode       string `json:"postalCode"`                 // 邮政编码，对应美国的zipcode。最大支持10个字符。美国zipcode支持标准5位格式（如94105）或带+4扩展代码的格式（如94105-1234）。当使用+4格式时，前5位用于业务校验和路线规划，后4位仅作为记录保存在订单数据中。SwiftX服务覆盖邮编，请联系商务获取。
}

// Validate 地址信息校验
func (m Address) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required.Error("姓名不能为空"), validation.Length(1, 100).Error("姓名长度不能超过 {{.max}} 个字符")),
		validation.Field(&m.PhoneNumber, validation.When(m.PhoneNumber != "", validation.Length(1, 50).Error("电话号码长度不能超过 {{.max}} 个字符"))),
		validation.Field(&m.RegionCode, validation.Required.Error("国家代码不能为空"), is.CountryCode2.ErrorObject(validation.NewError("422", "无效的国家代码 {{.value}}").SetParams(map[string]interface{}{"value": m.RegionCode}))),
		validation.Field(
			&m.StateProvince,
			validation.Required.Error("州/省不能为空"),
			validation.
				When(
					m.RegionCode == "US",
					validation.Match(
						regexp.MustCompile("^[A-Z]{2}$")).
						ErrorObject(
							validation.
								NewError("422", "无效的州/省 {{.value}}").
								SetParams(map[string]any{"value": m.StateProvince}),
						),
				).
				Else(
					validation.Length(1, 10).Error("州/省长度不能超过 {{.max}} 个字符"),
				),
		),
		validation.Field(&m.City, validation.Required.Error("城市不能为空"), validation.Length(1, 100).Error("城市长度不能超过 {{.max}} 个字符")),
		validation.Field(&m.District, validation.When(m.District != "", validation.Length(0, 100).Error("区域/县长度不能超过 {{.max}} 个字符"))),
		validation.Field(&m.StreetAddress, validation.Required.Error("街道地址不能为空"), validation.Length(0, 255).Error("街道地址长度不能超过 {{.max}} 个字符")),
		validation.Field(&m.Building, validation.When(m.Building != "", validation.Length(0, 255).Error("建筑物名称长度不能超过 {{.max}} 个字符"))),
		validation.Field(&m.PostalCode, validation.Required.Error("邮编不能为空"), validation.Length(1, 10).ErrorObject(validation.NewError("422", "邮编长度不能大于 {{.max}} 个字符"))),
	)
}
