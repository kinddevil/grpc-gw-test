package cluster

import (
	"context"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"strings"
	"time"

	etcv3 "go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

const (
	SCHEMA          = "GRPC_V3_LB"
	CLIENT_TIME_OUT = 5 // seconds
)

func NewBuilder(targets []string) resolver.Builder {
	if cli, err := initCli(targets); err != nil {
		panic(err)
	} else {
		return &etcdBuilder{targets: targets, client: cli}
	}
}

func initCli(targets []string) (*etcv3.Client, error) {
	return etcv3.New(etcv3.Config{
		Endpoints:   targets,
		DialTimeout: CLIENT_TIME_OUT * time.Second,
	})
}

type etcdBuilder struct {
	targets []string // register cluster address
	client  *etcv3.Client
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	if b.client == nil {
		var err error
		if b.client, err = initCli(b.targets); err != nil {
			panic(err) // cannot create the etcd connection at bootstrap
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	key := "/" + target.Scheme + "/" + target.Endpoint
	addrListStr, err := getStrList(key, b.client) // {'x.x.x.x:xxxx', ...}
	if err != nil {
		panic(err) // cannot get server list at bootstrap
	}

	addressList := string2Addr(addrListStr)

	d := &etcdResolver{
		ctx:    ctx,
		cancel: cancel,
		cc:     cc,
		client: b.client,
	}

	d.cc.UpdateState(resolver.State{Addresses: addressList})
	go d.watch(key, addressList)
	return d, nil
}

func getStrList(key string, cli *etcv3.Client) (targets []string, err error) {
	res, err := cli.Get(context.Background(), key, etcv3.WithPrefix())
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

func (b *etcdBuilder) Scheme() string { return SCHEMA }

type etcdResolver struct {
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	client *etcv3.Client
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOption) {}

func (r *etcdResolver) Close() {
	r.cancel()
	r.client.Close()
}

func (r *etcdResolver) watch(key string, addrList []resolver.Address) {

	rch := r.client.Watch(r.ctx, key, etcv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			// Cannot use string(ev.Kv.Value) for delete
			addr := strings.TrimPrefix(string(ev.Kv.Key), key)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			}
			//log.Printf("watch register servers %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			log.Printf("Update address list type %s key %q : %q to %v", ev.Type, ev.Kv.Key, ev.Kv.Value, addrList)
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
