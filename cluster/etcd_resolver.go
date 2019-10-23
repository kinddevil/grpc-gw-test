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

var (
	cli *etcv3.Client
)

func NewBuilder(targets []string) resolver.Builder {
	if err := initCli(targets); err != nil {
		panic(err)
	}
	return &etcdBuilder{targets: targets}
}

func initCli(targets []string) error {
	if cli == nil {
		var err error
		cli, err = etcv3.New(etcv3.Config{
			Endpoints:   targets,
			DialTimeout: 2 * time.Second, // TODO config dial timeout
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type etcdBuilder struct {
	targets []string // register cluster address
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	log.Println(target.Endpoint)
	log.Println(target.Scheme)

	if cli == nil {
		if err := initCli(b.targets); err != nil {
			panic(err) // cannot create the etcd connection at bootstrap
		}
	}

	key := target.Scheme + "/" + target.Endpoint
	addrListStr, err := getStrList(key) // {'x.x.x.x:xxxx', ...}
	if err != nil {
		panic(err) // cannot get server list at bootstrap
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := &etcdResolver{
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
	}


	d.cc.NewAddress(string2Addr(addrListStr))
	go d.watch(key)
	return d, nil
}

func getStrList(key string) (targets []string, err error) {
	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   targets,
		DialTimeout: 2 * time.Second,
	})

	res, err := cli.Get(context.TODO(), key, etcv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, v := range res.Kvs { // e.g. key:"foo" create_revision:16 mod_revision:55 version:2 value:"bar"
		targets = append(targets, string(v.Value))
	}
	return
}

func string2Addr(strList []string) (addrListObj []resolver.Address) {
	for _, addr := range strList {
		addrListObj = append(addrListObj, resolver.Address{Addr: addr})
	}
	return
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

func (r *etcdResolver) watch(key string) {

}
