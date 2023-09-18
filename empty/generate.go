package empty

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mempty.proto=github.com/mniak/duplicomp/empty --go-grpc_out=. empty.proto
