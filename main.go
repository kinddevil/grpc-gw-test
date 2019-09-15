package main

import (
	"flag"
	"grpc-gw-test/configs"
	"grpc-gw-test/servers"
)

func main() {
	env := flag.String("env", "dev", "environment: dev|staging|smoke|production|docker")
	flag.Parse()
	configs.CONFIGS = configs.LoadConfigs(env, "./resources")
	servers.StartServers(configs.CONFIGS)
}
