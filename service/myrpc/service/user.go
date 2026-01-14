package service

import (
	"context"

	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
	"github.com/entrehuihui/grpa-gateway-complete2/service/operate"
)

// User 列表
func (s Service) GetUser(ctx context.Context, in *proto.GetUserReq) (*proto.GetUserResp, error) {
	return operate.GetUser(ctx, in)
}

// User 添加
func (s Service) PostUser(ctx context.Context, in *proto.PostUserReq) (*proto.PostUserResp, error) {
	return operate.PostUser(ctx, in)
}

// User 修改
func (s Service) PutUser(ctx context.Context, in *proto.PutUserReq) (*proto.PutUserResp, error) {
	return operate.PutUser(ctx, in)
}
