package kafka

import (
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type protoEncoderImpl struct {
	msg proto.Message
}

func (p *protoEncoderImpl) Encode() ([]byte, error) {
	return proto.Marshal(p.msg)
}

func (p *protoEncoderImpl) Length() int {
	return proto.Size(p.msg)
}

// ProtoEncoder encodes Protobuf messages for Sarama clients
func ProtoEncoder(msg proto.Message) sarama.Encoder {
	return &protoEncoderImpl{msg: msg}
}
