package cluster

import (
	"context"

	"google.golang.org/grpc/resolver"
)

const (
	schema = "GRPC_ETCD"
)

func NewBuilder() resolver.Builder {
	return &etcdBuilder{}
}

type etcdBuilder struct {
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &etcdResolver{
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
	}
	d.cc.NewAddress([]resolver.Address{{Addr: "127.0.0.1" + ":" + "2378"}})
	return d, nil
}

func (b *etcdBuilder) Scheme() string { return schema }

type etcdResolver struct {
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOption) {

}

func (r *etcdResolver) Close() {

}
