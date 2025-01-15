package kafka

import (
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type protoDecoderImpl struct {
	msg proto.Message
}

func (p *protoDecoderImpl) Encode() ([]byte, error) {
	return proto.Marshal(p.msg)
}

func (p *protoDecoderImpl) Length() int {
	return proto.Size(p.msg)
}

// ProtoDecoder decodes Protobuf messages
func ProtoDecoder(msg proto.Message) sarama.Encoder {
	return &protoDecoderImpl{msg: msg}
}
