package servers

import (
	"context"
	"github.com/bsm/grpclb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/metadata"
	pb "grpc-gw-test/service_interfaces"
	"testing"
	"time"
)

func TestServeGRPC(t *testing.T) {
	defer setUp(t)(t)

	terminate := make(chan CancelFun, 1)
	defer close(terminate)

	go ServeGRPC(terminate, testCfg)

	testGrpcClient(t, testCfg.GetString("rest.grpc_addr"))

	terminateFunc := <-terminate
	terminateFunc()
}

func testGrpcClient(t *testing.T, address string) {
	// Set up a connection to the server.

	// Self defined balancer
	//localLB := grpclb.PickFirst(&grpclb.Options{
	//	Address: "127.0.0.1:8383",
	//})
	//balancer.Register(xxx)

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		//grpc.WithBalancerName(roundrobin.Name),
		//grpc.WithBalancer(localLB),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		t.Errorf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSampleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Meta-TrackId", "client-sent")

	_, err = c.Sample(ctx, &pb.Request{Id: "1", Name: "anonymous"})
	if err != nil {
		t.Errorf("could not score: %v", err)
	}
	// t.Logf("Test results: %v %v %v", r.Status, r.Code, r.Msg)
}
