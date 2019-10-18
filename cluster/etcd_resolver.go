package cluster

import (
	"context"
	"log"
	"time"

	etcv3 "go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

const (
	Schema = "/GRPC/ETCD"
)

func NewBuilder(targets []string) resolver.Builder {
	return &etcdBuilder{targets: targets}
}

type etcdBuilder struct {
	targets []string // register cluster address
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	log.Println("in etcd builder...")
	log.Println(target.Endpoint)
	log.Println(target.Scheme)
	log.Println(target.Authority)

	//configTargets := strings.Split(strings.ReplaceAll(target.Endpoint," ",""), ",")

	ctx, cancel := context.WithCancel(context.Background())
	d := &etcdResolver{
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
	}
	d.cc.NewAddress([]resolver.Address{{Addr: "127.0.0.1" + ":" + "2378"}})
	return d, nil
}

func updateTarget(targets []string, key string) {
	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   targets,
		DialTimeout: 2 * time.Second,
	})

	res, err := cli.Get(context.TODO(), key, etcv3.WithPrefix())
	if err != nil {
		//t.Fatal(err)
	}
	log.Println(res)
	//t.Log(res.Header)
	//t.Log(res.Count)
	//for _, v := range res.Kvs {
	//	t.Log(v)
	//}
}

func (b *etcdBuilder) Scheme() string { return Schema }

type etcdResolver struct {
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOption) {

}

func (r *etcdResolver) Close() {

}
