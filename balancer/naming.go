package balancer

import (
	"context"
	"fmt"
	etcv3 "go.etcd.io/etcd/clientv3"
	"log"
	"strings"
	"time"
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

	ticker := time.NewTicker(time.Second * time.Duration(ttl))

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
				log.Printf("Lease key already exisit %v", key)
			}

			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name string, addr string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	fmt.Printf("key:%v\n", "/"+SCHEMA+"/"+name+"/"+addr)
	_, err = cli.Put(context.Background(), "/"+SCHEMA+"/"+name+"/"+addr, addr, etcv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	resp, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	} else {
		<-resp
	}
	return nil
}

/**
UnRegister remove service from etcd
*/
func UnRegister(name string, addr string) {
	if cli != nil {
		key := "/"+SCHEMA+"/"+name+"/"+addr
		if _, err := cli.Delete(context.Background(), key); err != nil {
			log.Printf("Delete etcd server config error %v for key %v", err, key)
		}
	}
}
