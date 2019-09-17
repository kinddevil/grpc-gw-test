package servers

import (
	"errors"
	"github.com/spf13/viper"
	"grpc-gw-test/configs"
	"syscall"
	"testing"
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

func TestStartServers(t *testing.T) {

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	err := StartServers(testCfg, []*ServerInfo{
		{
			Name: "succ",
			Server: func(terminate chan<- func() error, cfg *viper.Viper) {
				terminate <- func() error { return nil }
			},
		},
	})
	if err != nil {
		t.Errorf("Should return nil after server starting - %v", err)
	}

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	err = StartServers(testCfg, []*ServerInfo{
		{
			Name: "fail",
			Server: func(terminate chan<- func() error, cfg *viper.Viper) {
				terminate <- func() error { return errors.New("terminate error test") }
			},
		},
	})
	if err == nil {
		t.Errorf("Should return non-empty after server starting")
	}
}
