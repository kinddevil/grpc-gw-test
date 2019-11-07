package servers

import (
	"errors"
	"github.com/spf13/viper"
	"grpc-gw-test/configs"
	"os"
	"syscall"
	"testing"
	"time"
)

var (
	testCfg *viper.Viper
)

func setUp(t *testing.T) func(t *testing.T) {
	env := "dev"
	testCfg = configs.LoadConfigs(&env, "../resources")

	return func(t *testing.T) {
		testCfg = nil
	}
}

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("signal was %v, want %v", s, sig)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for %v", sig)
	}
}

func TestStartServers(t *testing.T) {

	var err error
	go func() {
		err = StartServers(testCfg, []*ServerInfo{
			{
				Name: "succ",
				Server: func(terminate chan<- func() error, cfg *viper.Viper) {
					terminate <- func() error { return nil }
				},
			},
		})
	}()

	time.Sleep(10 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	time.Sleep(10 * time.Millisecond)

	if err != nil {
		t.Errorf("Should return nil after server starting - %v", err)
	}

	go func() {
		err = StartServers(testCfg, []*ServerInfo{
			{
				Name: "fail",
				Server: func(terminate chan<- func() error, cfg *viper.Viper) {
					terminate <- func() error { return errors.New("Terminate error test") }
				},
			},
		})
	}()

	time.Sleep(10 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	time.Sleep(10 * time.Millisecond)

	if err == nil {
		t.Errorf("Should return non-empty after server starting")
	}
}
