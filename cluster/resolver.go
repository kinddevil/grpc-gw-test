package cluster

import (
	"context"
	etcv3 "go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

//type Resolver interface {
//	// Resolve creates a Watcher for target.
//	Resolve(target string) (Watcher, error)
//}

type Resolver struct {}

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
