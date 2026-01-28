package service

import (
	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
	"github.com/entrehuihui/grpa-gateway-complete2/service/operate"
)

// 双向流式RPC
func (s Service) BiDiStream(stream proto.Stream_BiDiStreamServer) error {
	return operate.BiDiStream(stream)
}

// 客户端流式RPC
func (s Service) ClientStream(stream proto.Stream_ClientStreamServer) error {
	return operate.ClientStream(stream)
}

// 服务端流式RPC
func (s Service) UniStream(in *proto.StreamRequest, stream proto.Stream_UniStreamServer) error {
	return operate.UniStream(in, stream)
}
