package balancer

import (
	"google.golang.org/grpc/resolver"
	"testing"
)

func TestResolve(t *testing.T) {
	target := &resolver.Target{Scheme: SCHEMA, Endpoint: "TestEndpoint"}

	r := NewResolver("localhost:2378")
	if _, err := r.Build(*target, nil, resolver.BuildOption{}); err != nil {
		t.Errorf("TestResolve - Get error of build resolver %v", err)
	}

	//watch(r.(*IResolver), "test")

	//cli, err := etcv3.New(etcv3.Config{
	//	Endpoints:   []string{target.Endpoint},
	//	DialTimeout: 2 * time.Second,
	//})
	//
	//res, err := cli.Get(context.TODO(), "foo")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(res.Header)
	//t.Log(res.Count)
	//for _, v := range res.Kvs {
	//	t.Log(v)
	//}
}
