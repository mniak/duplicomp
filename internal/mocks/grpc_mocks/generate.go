package grpc_mocks

//go:generate protoc -I ../../cmd/sample/internal/ --go_out=. --go_opt=paths=source_relative --go_opt=Mping.proto=github.com/mniak/duplicomp/mocks/grpc_mocks ping.proto
