package service

import (
	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
	"github.com/entrehuihui/grpa-gateway-complete2/service/operate"
)

// SensorPost 流式测试
func (s Service) SensorPost(stream proto.Sensor_SensorPostServer) error {
	return operate.SensorPost(stream)
}
