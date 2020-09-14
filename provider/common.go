package provider

//go:generate protoc -I ../assets/grpc --go_out=plugins=grpc:. ../assets/grpc/provider.proto
