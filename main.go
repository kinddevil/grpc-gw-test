package main

import (
	"grpc_tpl/configs"
	"log"
)

func main() {
	// TODO load configs
	configs.CONFIGS.LoadConfigs()
	log.Println(configs.CONFIGS.Configs.GetString("grpc.port"))
	//servers.StartServers()
}
