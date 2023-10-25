package dynpb

import (
	_ "embed"
)

type Example struct {
	Bytes     []byte
	Data      []byte
	ProtoFile []byte
}

//go:embed examples/1/bytes.bin
var __example1_Bytes []byte

//go:embed examples/1/data.txt
var __example1_Data []byte

//go:embed examples/1/types.proto
var __example1_TypesProto []byte

var example1 = Example{
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

var example2 = Example{
	Bytes:     __example2_Bytes,
	Data:      __example2_Data,
	ProtoFile: __example2_TypesProto,
}
