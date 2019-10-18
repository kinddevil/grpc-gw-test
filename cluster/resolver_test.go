package cluster

import (
	"context"
	"google.golang.org/grpc/resolver"
	"testing"
	"time"

	etcv3 "go.etcd.io/etcd/clientv3"
)

func TestResolve(t *testing.T) {
	target := &resolver.Target{Scheme: SCHEMA, Endpoint: "127.0.0.1:2379"}
	//
	//r := NewResolver("localhost:2378")
	//if _, err := r.Build(*target, nil, resolver.BuildOption{}); err != nil {
	//	t.Errorf("TestResolve - Get error of build resolver %v", err)
	//}

	//watch(r.(*IResolver), "test")

	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   []string{target.Endpoint},
		DialTimeout: 2 * time.Second,
	})

	res, err := cli.Get(context.TODO(), "foo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Header) // e.g. cluster_id:14841639068965178418 member_id:10276657743932975437 revision:89 raft_term:7
	t.Log(res.Count)
	for _, v := range res.Kvs {
		t.Log(v) // e.g. key:"foo" create_revision:16 mod_revision:55 version:2 value:"bar"
	}
}
