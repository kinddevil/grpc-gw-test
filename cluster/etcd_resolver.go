package cluster

import (
	"google.golang.org/grpc/resolver"
)

func NewBuilder() resolver.Builder {
	return &etcdBuilder{}
}

type etcdBuilder struct {
}

func (r *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	return nil, nil
}

func (r *etcdBuilder) Scheme() string { return "" }
