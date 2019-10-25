package cluster

import (
	"context"
	etcv3 "go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"testing"
	"time"
)

type ClientConnMock struct {
	addresses     []resolver.Address
	serviceConfig string
}

func (c *ClientConnMock) UpdateState(state resolver.State) {
	c.addresses = state.Addresses
}

func (c *ClientConnMock) NewAddress(addresses []resolver.Address) {
	c.addresses = addresses
}

func (c *ClientConnMock) NewServiceConfig(serviceConfig string) {
	c.serviceConfig = serviceConfig
}

func TestResolve(t *testing.T) {
	etcdServer := "localhost:2379"
	//resolver.Register(NewBuilder([]string{etcdServer}))

	testBuilder := NewBuilder([]string{etcdServer})
	testEndpoint := "testEnd"
	target := &resolver.Target{Scheme: testBuilder.Scheme(), Endpoint: testEndpoint}

	clientConn := &ClientConnMock{}
	buildOpt := &resolver.BuildOption{}

	// Test resolver builder
	testResolver, err := testBuilder.Build(*target, clientConn, *buildOpt)
	defer testResolver.Close()

	if err != nil {
		t.Error("resolver builder should build a non-empty resolver...", err)
	}
	if len(clientConn.addresses) != 0 {
		t.Error("Should have no address info, get ", clientConn.addresses)
	}

	etcdClient, err := etcv3.New(etcv3.Config{
		Endpoints:   []string{etcdServer},
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		t.Error("Cannot connect to etcd...", err)
	}

	key := "/" + testBuilder.Scheme() + "/" + testEndpoint + "/"
	// Test adding server
	etcdClient.Put(context.TODO(), key, "someval")
	time.Sleep(500 * time.Microsecond)
	if len(clientConn.addresses) != 1 {
		t.Error("Should have only 1 address info, get ", clientConn.addresses)
	}

	// Test deleting server
	etcdClient.Delete(context.TODO(), key)
	time.Sleep(500 * time.Microsecond)
	if len(clientConn.addresses) != 0 {
		t.Error("Should have no address info after delete, get ", clientConn.addresses)
	}
}
