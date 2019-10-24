package cluster

import (
	"context"
	"fmt"
	etcv3 "go.etcd.io/etcd/clientv3"
	"log"
	"strings"
	"time"
)

const (
	SCHEMA          = "GRPC_V3_LB"
	CLIENT_TIME_OUT = 5 // seconds
)

/**
	Register register service with name as prefix to etcd, multi etcd addr should use ; to split
**/
func Register(etcdAddr, name, addr string, ttl int64) error {
	var err error

	if cli == nil {
		cli, err = etcv3.New(etcv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: CLIENT_TIME_OUT * time.Second,
		})
		if err != nil {
			return err
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl - 2))

	// TODO close at exist
	go func() {
		key := "/" + SCHEMA + "/" + name + "/" + addr
		for {
			getResp, err := cli.Get(context.Background(), key)
			if err != nil {
				log.Println(err)
			} else if getResp.Count == 0 {
				log.Printf("Init lease for key %v...", key)
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Println(err)
				}
			} else {
				//log.Printf("Lease key already exist %v", key)
			}

			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name string, addr string, ttl int64) error {
	ctx := context.Background()
	leaseResp, err := cli.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	log.Printf("key:%v\n", "/"+SCHEMA+"/"+name+"/"+addr)
	_, err = cli.Put(ctx, "/"+SCHEMA+"/"+name+"/"+addr, addr, etcv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	ch, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			ka := <-ch
			fmt.Println("ttl:", ka.TTL)
		}
	} ()

	return nil
}

/**
UnRegister remove service from etcd
*/
func UnRegister(name string, addr string) {
	if cli != nil {
		key := "/" + SCHEMA + "/" + name + "/" + addr
		if _, err := cli.Delete(context.Background(), key); err != nil {
			log.Printf("Delete etcd server config error %v for key %v", err, key)
		}
	}
}
