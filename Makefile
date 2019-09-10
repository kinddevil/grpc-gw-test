PROJECT=grpc_sample
GOPATH ?= $(shell go env GOPATH)

.PHONY: build

build:
	protoc -I sys_protos -I protos/ -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:service_interfaces --swagger_out=logtostderr=true:./specs --grpc-gateway_out=logtostderr=true:service_interfaces protos/sample.proto