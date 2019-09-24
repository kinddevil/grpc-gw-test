PROJECT=grpc_sample
GOPATH ?= $(shell go env GOPATH)

.PHONY: build

build:
	protoc -I sys_protos -I protos/ -I$(GOPATH)/src -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:service_interfaces --swagger_out=logtostderr=true:./specs --grpc-gateway_out=logtostderr=true:service_interfaces protos/sample.proto

test:
	docker images | grep "bitnami/etcd" || docker pull bitnami/etcd
	docker ps | grep "bitnami/etcd" || docker run -d --name etcd-server -p 2379:2379 -p 2380:2380 -e ALLOW_NONE_AUTHENTICATION=yes -e ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 bitnami/etcd
	mkdir -p reports && go test -cover -v -coverprofile ./reports/cover.out ./...
	go tool cover -html=./reports/cover.out -o ./reports/cover.html

valid:
	go vet -v ./...

fmt:
	go fmt ./...

bench:
	go test -bench=".*"