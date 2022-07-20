package serializer

import (
	"go-pcbook-micro/pb"
	"go-pcbook-micro/sample"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	// 写入 二进制文件
	laptop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	// 读取
	laptop2 := &pb.Laptop{}
	err = ReadProtobuffFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	// 写入 json 文件
	err = WriteProtobuffToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)
}
