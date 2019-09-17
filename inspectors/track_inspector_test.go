package inspectors

import (
	"context"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), map[string][]string{})
	if newCtx, err := GetUserInfo(ctx, nil); err != nil {
		t.Errorf("Track inspector error %v", err)
	} else {
		md, ok := metadata.FromIncomingContext(newCtx)
		if !ok {
			t.Errorf("Track inspector - Should have metadata in the context")
		}
		if _, ok := md["x-meta-trackid"]; !ok {
			t.Errorf("Track inspector - Should have x-meta-trackid in the metadata")
		}
	}
}
