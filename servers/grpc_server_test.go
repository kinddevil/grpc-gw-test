package servers

import (
	"log"
	"testing"
)

func TestServeGRPC(t *testing.T) {
	terminate := make(chan<- CancelFun, 1)
	defer close(terminate)

	log.Println(terminate)
	//go ServeGRPC(terminate)
}
