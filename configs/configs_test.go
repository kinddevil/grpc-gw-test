package configs

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestLoadConfigsSucc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	env := "dev"
	cfg := LoadConfigs(&env, "../resources")
	if cfg.Get("grpc.port") != ":50051" {
		t.Errorf("Setting config error - grpc.port should equal %v, get %v", ":50051", cfg.Get("grpc.port"))
	}

	// Test file not found
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should give file not found error")
		}
	}()
	LoadConfigs(&env,"whatever")

	// Test read config error
	oldLoadConfigs := loadConfigs
	//defer func () {
	//	loadConfigs = oldLoadConfigs
	//}()
	loadConfigs := func(path, filename string, defaults map[string]interface{}) (*viper.Viper, error) {
		return nil, errors.New("Test read config with return error")
	}
	LoadConfigs(&env,"../resources")
}
