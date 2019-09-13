package services

import (
	"context"
	pb "grpc-gw-test/service_interfaces"
	"testing"
)

func SampleTest(t *testing.T) {
	s := Server{}

	data := []struct {
		id     string
		name   string
		status int32
	}{
		{
			id:     "1",
			name:   "Hello world",
			status: 1,
		},
		{
			id:     "2",
			name:   "你好 世界",
			status: 1,
		},
	}

	for _, item := range data {
		req := &pb.Request{Id: item.id, Name: item.name}
		resp, err := s.Sample(context.Background(), req)
		if err != nil {
			t.Errorf("SampleTest(%v) got unexpected error", item.id)
		}
		if resp.Status != item.status {
			t.Errorf("SampleTest(%v)=%v, wanted %v", item.status, resp.Status, item.id)
		}
	}
}
