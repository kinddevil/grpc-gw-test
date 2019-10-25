package cluster

import (
	"context"
	"errors"
	etcv3 "go.etcd.io/etcd/clientv3"
	"log"
	"strings"
	"time"
)

type RegisterService struct {
	client *etcv3.Client
	lease  *etcv3.LeaseKeepAliveResponse
}

/**
	Register register service with name as prefix to etcd, multi etcd addr should use ; to split
**/
func (s *RegisterService) Register(etcdAddr, name, addr string, ttl int64) error {
	var err error

	// TODO not use global client
	if s.client == nil {
		s.client, err = etcv3.New(etcv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: CLIENT_TIME_OUT * time.Second,
		})
		if err != nil {
			return err
		}
	}

	key := "/" + SCHEMA + "/" + name + "/" + addr
	getResp, err := s.client.Get(context.Background(), key)

	if err != nil {
		return err
	} else if getResp.Count == 0 {
		log.Printf("Init lease for key %v...", key)
		leaseChan, err := s.withAlive(name, addr, ttl)
		if err != nil {
			return err
		}
		go s.monitorLease(leaseChan)
	} else {
		// duplicated service
		return errors.New("Service already exist " + key)
	}

	return nil
}

func (s *RegisterService) withAlive(name string, addr string, ttl int64) (<-chan *etcv3.LeaseKeepAliveResponse, error) {
	ctx := context.Background()
	leaseResp, err := s.client.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}

	log.Printf("register for key:%v\n", "/"+SCHEMA+"/"+name+"/"+addr)
	_, err = s.client.Put(ctx, "/"+SCHEMA+"/"+name+"/"+addr, addr, etcv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	ch, err := s.client.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (s *RegisterService) monitorLease(leaseChan <-chan *etcv3.LeaseKeepAliveResponse) {
	// Fix bug of etcd queue is full
	for {
		if lease, ok := <-leaseChan; !ok {
			break
		} else {
			s.lease = lease
		}
	}
}

/**
UnRegister remove service from etcd
*/
func (s *RegisterService) UnRegister(name string, addr string) {
	if s.client != nil {
		ctx := context.Background()

		if s.lease != nil {
			s.client.Revoke(ctx, s.lease.ID)
		}
		key := "/" + SCHEMA + "/" + name + "/" + addr
		if _, err := s.client.Delete(ctx, key); err != nil {
			log.Printf("Delete etcd server config error %v for key %v", err, key)
		}
	}
}
