package grpc

/* Install protoc and
sudo pacman -S protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
*/

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go_opt=Mping.proto=github.com/mniak/duplicomp/internal/samples/grpc --go-grpc_out=. ping.proto

// //go:generate mockgen -package=grpc -destination=mock_pinger_server_test.go . PingerServer
// //go:generate mockgen -package=grpc -destination=mock_pinger_client_test.go . PingerClient
