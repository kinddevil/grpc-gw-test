package cluster

import (
	"context"
	"fmt"
	"testing"
	"time"

	etcv3 "go.etcd.io/etcd/clientv3"
)

func TestNaming(t *testing.T) {
	target := "http://localhost:2379"
	name := "test"
	addr := "address"
	key := fmt.Sprintf("/%v/%v/%v", SCHEMA, name, addr)

	rServer := &RegisterService{}
	rServer.Register(target, "test", "address", 5)
	time.Sleep(time.Second)

	cli, err := etcv3.New(etcv3.Config{
		Endpoints:   []string{target},
		DialTimeout: 1 * time.Second,
	})

	res, err := cli.Get(context.TODO(), key, etcv3.WithPrefix())
	if err != nil {
		t.Fatal(err)
	}

	if res.Count != 1 {
		t.Errorf("TestNaming - Shouldl only have one key after register %v", key)
	}
	t.Logf("TestNaming - headers, count, more - %v - %v - %v...", res.Header, res.Count, res.More)
	for _, v := range res.Kvs {
		t.Logf("TestNaming - value - %v", v)
	}

	rServer.UnRegister(name, addr)

	res, err = cli.Get(context.TODO(), key, etcv3.WithPrefix())
	if err != nil {
		t.Fatal(err)
	}

	if res.Count != 0 {
		t.Errorf("TestNaming - Shouldl have no key after unregister %v", key)
	}
}
