package cluster

import (
	"context"
	etcv3 "go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

//type Resolver interface {
//	// Resolve creates a Watcher for target.
//	Resolve(target string) (Watcher, error)
//}

type Resolver struct {
	rawAddr string
	cc      resolver.ClientConn
}

type ClientConn struct {}

func NewResolver(addr string) resolver.Builder {
	return Resolver{rawAddr: addr}
}

func (r *Resolver) UpdateState(resolver.State) {}

func (r *Resolver) NewAddress(addresses []resolver.Address) {}

func (r *Resolver) NewServiceConfig(serviceConfig string) {}

func (r *Resolver) ResolveNow(resolver.ResolveNowOption) {}

func (r *Resolver) Close() {}

func (r *Resolver) Build(target resolver.Target, cc ClientConn, opts resolver.BuildOption) (Resolver, error) {
	return Resolver{}, nil
}

func (r *Resolver) Scheme() string { return "" }

func (r *Resolver) Resolve(target string) {

	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   []string{target},
		DialTimeout: 2 * time.Second,
	})

	if err == nil {
		// handle errors
	}

	defer cli.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	_, err = cli.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
}
