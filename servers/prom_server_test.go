package servers

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServerProm(t *testing.T) {
	defer setUp(t)(t)

	terminate := make(chan CancelFun, 1)
	defer close(terminate)

	addr := fmt.Sprintf("http://localhost%v", testCfg.GetString("common.prom.port"))

	go ServeProm(terminate, testCfg)

	if _, err := http.Get(addr); err != nil {
		t.Errorf("Test http proxy error - %v", err)
	}

	terminateFunc := <-terminate
	terminateFunc()
}
