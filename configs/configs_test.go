package configs

import (
	"testing"
)

func TestLoadConfigsSucc(t *testing.T) {
	env := "dev"
	cfg := LoadConfigs(&env, "../resources")
	if cfg.Get("grpc.port") != ":50051" {
		t.Errorf("Setting config error - grpc.port should equal %v, get %v", ":50051", cfg.Get("grpc.port"))
	}
}

func TestLoadConfigsEnvErr(t *testing.T) {
	invalidEnv := "whatever_env"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should give env is not valid")
		}
	}()
	LoadConfigs(&invalidEnv, "../resources")
}

func TestLoadConfigsResourceErr(t *testing.T) {
	env := "dev"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should give file not found error")
		}
	}()
	LoadConfigs(&env, "whatever")
}
