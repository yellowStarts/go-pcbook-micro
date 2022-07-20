package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// ProtobuffToJSON protobuff 转化成 json
func ProtobuffToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true, // 大写字母驼峰是否转换下划线连接
	}

	return marshaler.MarshalToString(message)
}
