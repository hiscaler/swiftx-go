package swiftx

import (
	"context"
	"strconv"
)

// Ping 服务
type pingService service

// Pong 返回请求的数值，可用于请求测试/健康检查
func (s pingService) Pong(ctx context.Context, i int) (int, error) {
	var res int
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetQueryParam("i", strconv.Itoa(i)).
		SetResult(&res).
		Get("/pingPong")
	if err = recheckError(resp, err); err != nil {
		return 0, err
	}
	return res, nil
}
