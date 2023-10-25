#!/bin/bash
cd $(dirname $0)
protoc --encode=Object "types.proto" < "data.txt" > bytes.bin