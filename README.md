# gRPC Report Generation Service

## Requirements
- Go 1.20+
- `protoc` installed
- `protoc-gen-go` and `protoc-gen-go-grpc` installed
cmd - go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


## Run Instructions
1. Install Go modules: `go mod tidy`
2. Compile the proto: `protoc --go_out=. --go-grpc_out=. proto/*.proto`
3. Run the server: `go run server/*.go`
4. Test using `grpcurl`: 