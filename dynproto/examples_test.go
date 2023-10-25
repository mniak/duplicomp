package dynproto

import (
	_ "embed"

	"github.com/mniak/duplicomp/dynproto/examples"
)

//go:embed examples/1/bytes.bin
var __example1_Bytes []byte

//go:embed examples/1/data.txt
var __example1_Data []byte

//go:embed examples/1/types.proto
var __example1_TypesProto []byte

var example1 = examples.Example{
	Bytes:     __example1_Bytes,
	Data:      __example1_Data,
	ProtoFile: __example1_TypesProto,
}

//go:embed examples/2/bytes.bin
var __example2_Bytes []byte

//go:embed examples/2/data.txt
var __example2_Data []byte

//go:embed examples/2/types.proto
var __example2_TypesProto []byte

var example2 = examples.Example{
	Bytes:     __example2_Bytes,
	Data:      __example2_Data,
	ProtoFile: __example2_TypesProto,
}
