package configs

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

var (
	DEFAULT_CONFIGS = map[string]interface{}{
		"grpc": map[string]interface{}{
			"port":                ":50051",
			"max_connection_idle": "300",
		},
		"rest": map[string]interface{}{
			"port":      ":8081",
			"grpc_addr": "localhost:50051",
		},
	}

	CONFIGS = LoadConfigs()
)

func LoadConfigs() *viper.Viper {
	env := flag.String("env", "dev", "environment: dev|staging|smoke|production|docker")
	if !checkEnv(env) {
		panic("Invalid environment")
	}
	if config, err := loadConfigs("./resources", *env, DEFAULT_CONFIGS); err != nil {
		log.Fatal(err)
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
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
