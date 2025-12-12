package entity

// Address 收发货地址
type Address struct {
	Name             string `json:"name"`             // 姓名
	PhoneCountryCode int    `json:"phoneCountryCode"` // 电话国际区号。如果留空，将使用regionCode对应的国际区号作为电话的区号
	PhoneNumber      string `json:"phoneNumber"`      // 电话号码（寄件人地址可选，收件人地址必填）
	PhoneExtension   string `json:"phoneExtension"`   // 电话分机号
	Building         string `json:"building"`         // 建筑物名称，对应美国地址的address line 2
	StreetAddress    string `json:"streetAddress"`    // 街道地址，对应美国地址的address line 1
	District         string `json:"district"`         // 区域/县，对应美国地址的county
	City             string `json:"city"`             // 城市
	StateProvince    string `json:"stateProvince"`    // 州/省（美国州请使用2字母缩写，如CA、NY）
	PostalCode       string `json:"postalCode"`       // 邮政编码，对应美国的zipcode。最大支持10个字符。美国zipcode支持标准5位格式（如94105）或带+4扩展代码的格式（如94105-1234）。当使用+4格式时，前5位用于业务校验和路线规划，后4位仅作为记录保存在订单数据中。SwiftX服务覆盖邮编，请联系商务获取。
	RegionCode       string `json:"regionCode"`       // 国家编码，类似US、CN（Enum: "US" "CA" "CN"）
}
