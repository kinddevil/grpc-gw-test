package cluster

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"
	"testing"
	"time"

	pb "grpc-gw-test/service_interfaces"
)

func TestResolve(t *testing.T) {

	resolver.Register(NewBuilder([]string{"localhost:2379"}))

	resolver.SetDefaultScheme(SCHEMA)

	cc, err := grpc.Dial(
		SCHEMA+":///grpc-gw/",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	//grpc.WithBalancer(grpc.RoundRobin(r)),
	//grpc.WithBlock(),
	)
	defer cc.Close()
	log.Printf("test reolve connection %v, error %v", cc, err)

	ctx := context.TODO()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	for i := 0; i < 1002; i++ {
		c := pb.NewSampleServiceClient(cc)

		ret, err := c.Sample(ctx, &pb.Request{Id: "5", Name: "anonymous"})
		log.Printf("get ret %v with error %v", ret, err)
		time.Sleep(2 * time.Second)
	}
}
