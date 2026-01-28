package swiftx

import (
	"testing"

	"github.com/hiscaler/swiftx-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestOrderService_Create(t *testing.T) {
	// 寄件人地址信息
	senderAddress := SenderAddress{
		RegionCode:    "US",
		StateProvince: "CA",
		City:          "Ontario",
		StreetAddress: "2078 E Francis Street",
		PostalCode:    "91761",
		Name:          "ZEB2",
		PhoneNumber:   "1096398373",
	}
	// 收件人地址信息
	recipientAddress := RecipientAddress{
		RegionCode:    "US",
		StateProvince: "TX",
		City:          "Fort Worth",
		StreetAddress: "W1302 WELCH RD",
		PostalCode:    "76118",
		Name:          "Dak Jaech",
		PhoneNumber:   "+1 347-447-3197 ext. 07360",
	}
	// 创建订单请求数据
	req := CreateOrderRequest{
		OrderScope:        entity.OrderScopeDomestic,        // 订单类型：国内
		ServiceType:       entity.ServiceTypeExp,            // 服务类型：标快
		DeliveryMethod:    entity.DeliveryMethodHdy,         // 送货方式：上门派送
		CooperationMethod: entity.CooperationMethodMerchant, // 合作方式：商家
		InsuranceService: &InsuranceService{
			IsInsured: false,
			InsuredValue: Value{
				Amount:       100,
				CurrencyCode: "USD",
			},
		},
		PackageInfo: CreateOrderPackageInformation{ // 包裹信息
			SenderAddress:    senderAddress,    // 寄件人地址
			RecipientAddress: recipientAddress, // 收件人地址
			UseImperialUnit:  false,            // 不使用英制单位
			Weight:           1.5,              // 重量
			Length:           10,               // 长度
			Width:            10,               // 宽度
			Height:           5,                // 高度
			Value: Value{ // 费用金额
				Amount:       100,
				CurrencyCode: "USD",
			},
			CustomerName: "测试客户", // 客户名
			StoreName:    "测试店铺", // 店铺名
			SkuList: []CreateOrderPackageGoods{ // SKU 列表
				{
					Name:     "测试 SKU 1",
					Quantity: 1,
					Code:     "SKU001",
					Value:    null.FloatFrom(50),
				},
				{
					Name:     "测试 SKU 2",
					Quantity: 2,
					Code:     "SKU002",
					Value:    null.FloatFrom(25),
				},
			},
		},
		ShippingLabelInfo: ShippingLabelInformation{ // 运单印刷数据
			OrderNumber: "TEST-ORDER-12345", // 上游订单号
		},
	}
	order, err := client.Services.Order.Create(ctx, req)
	if err != nil {
		t.Logf("client.Services.Order.Create() 错误: %v", err) // 打印错误信息
		t.Fatalf("client.Services.Order.Create() 失败: %v", err)
	} else {
		assert.Equal(t, req.ShippingLabelInfo.OrderNumber, order.CustomerOrderNumber)
		assert.NotEmpty(t, order.ShipmentNumber)
		assert.NotEmpty(t, order.ShippingLabel)
	}
}

func TestOrderService_Cancel(t *testing.T) {
	shipmentNumber := "SWX475440000011278280"
	success, err := client.Services.Order.Cancel(ctx, shipmentNumber)
	if err != nil {
		t.Fatalf("client.Services.Order.Cancel() 错误: %v", err)
	}
	if !success {
		t.Error("期望成功取消订单，但操作失败")
	}
}

func TestOrderService_Tracking(t *testing.T) {
	// 请替换为测试环境中的有效运单号
	shipmentNumber := "SWX852250000011278331"
	results, err := client.Services.Order.Tracking(ctx, shipmentNumber)
	if err != nil {
		t.Fatalf("client.Services.Order.Tracking() 错误: %v", err)
	}
	if len(results) == 0 {
		t.Error("期望获取到物流跟踪结果，但结果为空")
	}
}

func TestOrderService_Price(t *testing.T) {
	// 请替换为测试环境中的有效运单号
	shipmentNumber := "SWX852250000011278331"
	results, err := client.Services.Order.Price(ctx, shipmentNumber)
	if err != nil {
		t.Fatalf("client.Services.Order.Price() 错误: %v", err)
	}
	if len(results) == 0 {
		t.Error("期望获取到订单价格，但结果为空")
	}
}
