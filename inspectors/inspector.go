package inspectors

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

type Inspector = func(ctx context.Context, in interface{}) (context.Context, error)

func MiddlewareFunc(middlewareFunc Inspector) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, in interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx, err := middlewareFunc(ctx, in)
		if err != nil {
			log.Printf("Inspecting func \"%v\" error %v", info.FullMethod, err)
			return handler(ctx, in)
		}
		return handler(newCtx, in)
	}
}