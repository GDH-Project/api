package grpc

// 	protoc        v6.32.0
// 	protoc-gen-go v1.36.10
//  protoc-gen-go-grpc v1.5.1

//go:generate protoc --proto_path=../../internal/api_spec/auth/proto --go_out=. --go-grpc_out=. ../../internal/api_spec/auth/proto/user.proto
//go:generate protoc --proto_path=../../internal/api_spec/auth/proto --go_out=. --go-grpc_out=. ../../internal/api_spec/auth/proto/auth.proto
