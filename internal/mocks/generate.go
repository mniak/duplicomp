package mocks

//go:generate protoc -I ../samples/internal/ --go_out=grpc_mocks --go_opt=paths=source_relative --go_opt=Mping.proto=github.com/mniak/duplicomp/mocks/grpc_mocks ping.proto
