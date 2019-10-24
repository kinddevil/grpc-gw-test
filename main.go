package main

import (
	"flag"
	"google.golang.org/grpc/resolver"
	"grpc-gw-test/cluster"
	"grpc-gw-test/configs"
	"grpc-gw-test/servers"
)

func main() {

	resolver.Register(cluster.NewBuilder([]string{"localhost:2379"}))

	resolver.SetDefaultScheme(cluster.SCHEMA)

	env := flag.String("env", "dev", "environment: dev|staging|smoke|production|docker")
	flag.Parse()
	configs.CONFIGS = configs.LoadConfigs(env, "./resources")
	servers.StartServers(configs.CONFIGS, servers.ServerList)
}
