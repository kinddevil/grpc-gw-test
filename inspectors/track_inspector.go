package inspectors

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func GetUserInfo(ctx context.Context, in interface{}) (newCtx context.Context, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		trackHeader := "X-Meta-TrackId"
		if len(md.Get(trackHeader)) == 0 {
			md.Append(trackHeader, uuid.Must(uuid.NewRandom()).String())
		}
	}
	newCtx = ctx
	return
}
