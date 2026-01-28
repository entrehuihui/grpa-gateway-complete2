package operate

import (
	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
)

// 双向流式RPC
func BiDiStream(stream proto.Stream_BiDiStreamServer) error {
	var err error
	return err
}

// 客户端流式RPC
func ClientStream(stream proto.Stream_ClientStreamServer) error {
	var err error
	return err
}

// 服务端流式RPC
func UniStream(in *proto.StreamRequest, stream proto.Stream_UniStreamServer) error {
	var err error
	return err
}
