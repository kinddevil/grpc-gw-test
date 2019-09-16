package servers

import (
	"github.com/spf13/viper"
	"grpc-gw-test/configs"
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

}
