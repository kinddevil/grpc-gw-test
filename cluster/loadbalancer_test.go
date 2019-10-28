package cluster

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	//"google.golang.org/grpc/metadata"
	"testing"
	//"time"

	//"context"
	//pb "grpc-gw-test/service_interfaces"

	"google.golang.org/grpc/resolver"
)

func TestLoadBalance(t *testing.T) {
	resolver.Register(NewBuilder([]string{"localhost:2379"}))

	conn, err := grpc.Dial("GRPC_V3_LB:///grpc-gw",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		t.Errorf("Did not connect: %v", err)
	}
	defer conn.Close()

	/*
		c := pb.NewSampleServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ctx = metadata.AppendToOutgoingContext(ctx, "X-Meta-TrackId", "client-sent")

		for i := 0; i < 2; i++ {
			_, err = c.Sample(ctx, &pb.Request{Id: "1", Name: "anonymous"})
			if err != nil {
				t.Errorf("could not score: %v", err)
			}
			// t.Logf("Test results: %v %v %v", r.Status, r.Code, r.Msg)
		}
	*/
}
