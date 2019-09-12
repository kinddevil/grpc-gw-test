package main

import (
	"grpc_tpl/configs"
	"grpc_tpl/servers"
	"log"
)

func main() {
	// TODO load configs
	configs.LoadConfigs()
	log.Println(configs.CONFIGS.GetString("grpc.port"))
	servers.StartServers()
}
