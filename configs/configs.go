package configs

import (
	"github.com/spf13/viper"
	"log"
)

var (
	DEFAULT_CONFIGS = map[string]interface{}{
		"grpc": map[string]interface{}{
			"port":                ":50051",
			"max_connection_idle": 300,
			"time_out":            7,
		},
		"rest": map[string]interface{}{
			"port":      ":8081",
			"grpc_addr": "localhost:50051",
		},
	}

	CONFIGS *viper.Viper
)

func LoadConfigs(env *string, configDir string) *viper.Viper {
	if !checkEnv(env) {
		panic("Invalid environment")
	}
	log.Printf("Init environments %v...", *env)

	if config, err := loadConfigs(configDir, *env, DEFAULT_CONFIGS); err != nil {
		panic(err)
	} else {
		return config
	}
	return nil
}

func checkEnv(env *string) bool {
	switch *env {
	case
		"dev",
		"staging",
		"smoke",
		"production",
		"docker":
		return true
	}
	return false
}

func loadConfigs(path, filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(path)
	v.AutomaticEnv() // Use sys environment as default
	err := v.ReadInConfig()
	return v, err
}
