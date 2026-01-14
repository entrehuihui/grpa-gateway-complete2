package service

import (
	"context"
	"grpa-gateway-complete2/service/myrpc/proto"
	"grpa-gateway-complete2/service/operate"
)

// GetRole 列表
func (s Service) GetRole(ctx context.Context, in *proto.GetRoleReq) (*proto.GetRoleResp, error) {
	return operate.GetRole(ctx, in)
}

// GetRole 添加
func (s Service) PostRole(ctx context.Context, in *proto.PostRoleReq) (*proto.PostRoleResp, error) {
	return operate.PostRole(ctx, in)
}

// GetRole 修改
func (s Service) PutRole(ctx context.Context, in *proto.PutRoleReq) (*proto.PutRoleResp, error) {
	return operate.PutRole(ctx, in)
}

// GetRole 列表
func (s Service) GetRoleAuth(ctx context.Context, in *proto.GetRoleAuthReq) (*proto.GetRoleAuthResp, error) {
	return operate.GetRoleAuth(ctx, in)
}
