package entity

const (
	Prod = "prod" // 生产环境
	Test = "test" // 测试环境
	Dev  = "dev"  // 开发环境
)

const (
	OrderScopeDomestic      = "DOMESTIC"      // 订单类型：国内
	OrderScopeInternational = "INTERNATIONAL" // 订单类型：国际
)

const (
	ServiceTypeEco = "ECO" // 服务类型：特惠
	ServiceTypeExp = "EXP" // 服务类型：标快
)

const (
	DeliveryMethodHdy = "HDY" // 送货方式：上门派送
	DeliveryMethodSpu = "SPU" // 送货方式：自提
)

const (
	CooperationMethodPlatform    = "PLATFORM"     // 合作方式：平台
	CooperationMethodMerchant    = "MERCHANT"     // 合作方式：商家
	CooperationMethodWesternPost = "WESTERN_POST" // 合作方式：西邮
)
