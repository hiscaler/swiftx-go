package entity

import (
	"github.com/hiscaler/swiftx-go/response"
)

// TrackingResult 物流跟踪结果
type TrackingResult struct {
	Result            response.Result
	TrackingNo        string  `json:"trackingNo"`        // 跟踪号
	TrackingEventList []Track `json:"trackingEventList"` // 跟踪信息
}
