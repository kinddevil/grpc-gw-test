package inspectors

import (
	"context"
	"testing"
)

func TestMiddlewareFunc(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error){return nil, nil}
	MiddlewareFunc(func(ctx context.Context, in interface{}) (context.Context, error){
		return nil, nil
	})(context.Background(), nil, nil, handler)
}
