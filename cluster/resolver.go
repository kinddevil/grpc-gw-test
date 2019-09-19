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

type ClientConn struct {}

func (cc *ClientConn) UpdateState(resolver.State) {

}

func (cc *ClientConn) NewAddress(addresses []resolver.Address) {

}

func (cc *ClientConn) NewServiceConfig(serviceConfig string) {}

type Resolver struct {

}

func (r *Resolver) ResolveNow(resolver.ResolveNowOption) {}

func (r *Resolver) Close() {}

type Builder struct {

}

func (b *Builder) Build(target resolver.Target, cc ClientConn, opts resolver.BuildOption) (Resolver, error) {
	return Resolver{}, nil
}

func (b *Builder) Scheme() string { return "" }

func (resolver *Resolver) Resolve(target string) {

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
