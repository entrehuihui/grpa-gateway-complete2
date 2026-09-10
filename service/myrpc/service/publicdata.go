package service

import (
	"context"

	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
	"github.com/entrehuihui/grpa-gateway-complete2/service/operate"
)

// GetCacheURL 获取临时上传连接
func (s Service) GetCacheURL(ctx context.Context, in *proto.GetCacheURLReq) (*proto.GetCacheURLResp, error) {
	return operate.GetCacheURL(ctx, in)
}
