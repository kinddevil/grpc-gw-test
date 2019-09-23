package cluster

import (
	"context"
	etcv3 "go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

func TestResolve(t *testing.T) {
	target := "http://localhost:2379"

	r := IResolver{}
	r.Resolve(target)

	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   []string{target},
		DialTimeout: 2 * time.Second,
	})

	res, err := cli.Get(context.TODO(), "foo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Header)
	t.Log(res.Count)
	for _, v := range res.Kvs {
		t.Log(v)
	}
}
