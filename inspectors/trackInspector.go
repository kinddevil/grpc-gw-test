package inspectors

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func GetUserInfo(ctx context.Context, in interface{}) (newCtx context.Context, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		md.Append("X-Meta-TrackId", "track")
		newCtx = context.WithValue(ctx, "X-Meta-TrackId", "track")
	} else {
		newCtx = ctx
	}
	return
}