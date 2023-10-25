//go:build mage

package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/samber/lo"
)

func GenerateExamples() {
	dir := lo.Must(os.ReadDir("."))
	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		dataTxt := lo.Must(os.Open(filepath.Join(entry.Name(), "data.txt")))
		defer dataTxt.Close()

		bytesBin := lo.Must(os.Create(filepath.Join(entry.Name(), "bytes.bin")))
		defer bytesBin.Close()

		cmdProtoc := exec.Command("protoc", "--encode=Object", filepath.Join(entry.Name(), "types.proto"))
		cmdProtoc.Stdin = dataTxt
		cmdProtoc.Stdout = bytesBin
		cmdProtoc.Stderr = os.Stderr

		lo.Must0(cmdProtoc.Run())
	}
}
