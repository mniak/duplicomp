package dynpb

import (
	"embed"
	_ "embed"
	"io"
	"path/filepath"

	"github.com/samber/lo"
)

type Example struct {
	Bytes     []byte
	Data      []byte
	ProtoFile []byte
}

//go:embed examples/*/*
var __examplesFS embed.FS

func LoadExample(name string) Example {
	bytesFile := lo.Must(__examplesFS.Open(filepath.Join("examples", name, "bytes.bin")))
	bytes := lo.Must(io.ReadAll(bytesFile))

	dataFile := lo.Must(__examplesFS.Open(filepath.Join("examples", name, "data.txt")))
	data := lo.Must(io.ReadAll(dataFile))

	protoFile := lo.Must(__examplesFS.Open(filepath.Join("examples", name, "types.proto")))
	proto := lo.Must(io.ReadAll(protoFile))

	return Example{
		Bytes:     bytes,
		Data:      data,
		ProtoFile: proto,
	}
}

func AllExamples() []Example {
	dir := lo.Must(__examplesFS.ReadDir("examples"))
	var result []Example
	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		result = append(result, LoadExample(filepath.Base(entry.Name())))
	}
	return result
}
