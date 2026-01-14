package service

import "github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"

// Service .
type Service struct {
	// ##继承
	proto.UnimplementedRoleServer
	proto.UnimplementedUserServer
	// ##继承
}

// NewService .
func NewService() *Service {
	s := new(Service)
	return s
}
	