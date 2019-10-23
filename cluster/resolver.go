package cluster

import (
	"context"
	etcv3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/naming"
	"google.golang.org/grpc/resolver"
	"log"
	"strings"
	"time"
)

const (
	CLIENT_TIME_OUT = 5 // seconds
)

type BaseResolver interface {
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error)
	Scheme() string
	ResolveNow(resolver.ResolveNowOption)
	Resolve(target string) (naming.Watcher, error)
	Close()
}

type CBuilder struct {
	RawAddr string
	cc      resolver.ClientConn
	state   resolver.State
}

func NewResolver(addr string) resolver.Builder {
	return &CBuilder{RawAddr: addr}
}

func (r *CBuilder) UpdateState(resolver.State) {}

func (r *CBuilder) NewAddress(addresses []resolver.Address) {}

func (r *CBuilder) NewServiceConfig(serviceConfig string) {}

func (r *CBuilder) ResolveNow(resolver.ResolveNowOption) {}

func (r *CBuilder) Close() {}

func (r *CBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	if cli == nil {
		cli, err = etcv3.New(etcv3.Config{
			Endpoints:   strings.Split(r.RawAddr, ";"),
			DialTimeout: CLIENT_TIME_OUT * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	r.cc = cc

	go watch(r, "/"+target.Scheme+"/"+target.Endpoint+"/")

	return r, nil
}

func watch(r *CBuilder, keyPrefix string) {
	addrList := make([]resolver.Address, 0, 1)

	log.Printf("get key from resolver %v", keyPrefix)
	getResp, err := cli.Get(context.Background(), keyPrefix, etcv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}

	r.cc.NewAddress(addrList)
	//r.cc.UpdateState(resolver.State{Addresses: addrList})

	rch := cli.Watch(context.Background(), keyPrefix, etcv3.WithPrefix())
	// TODO close when existing
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.NewAddress(addrList)
					//r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.NewAddress(addrList)
					//r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			}
			log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
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

func (r *CBuilder) Scheme() string { return SCHEMA }

func (r *CBuilder) Resolve(target string) (naming.Watcher, error) { return nil, nil }
