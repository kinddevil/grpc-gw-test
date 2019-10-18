package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"

	"grpc-gw-test/cluster"
)

func main() {
	//env := flag.String("env", "dev", "environment: dev|staging|smoke|production|docker")
	//flag.Parse()
	//configs.CONFIGS = configs.LoadConfigs(env, "./resources")
	//servers.StartServers(configs.CONFIGS, servers.ServerList)

	resolver.Register(cluster.NewBuilder([]string{"localhost:2379"}))

	//resolver.SetDefaultScheme("dns")
	resolver.SetDefaultScheme(cluster.Schema)

	//r := &cluster.EtcdResolver{}

	cc, err := grpc.Dial(
		cluster.Schema + ":///sample_server",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		//grpc.WithBalancer(grpc.RoundRobin(r)),
		//grpc.WithBlock(),
		)
	log.Println(cc)
	log.Println(err)
}
