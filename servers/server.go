package servers

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type CancelFun = func() error
type Server func(chan<- func() error, *viper.Viper)
type ServerInfo struct {
	name   string
	server Server
	cancel CancelFun
}

var(
	ServerList = []ServerInfo{
		//{
		//	name:   "grpc",
		//	server: ServeGRPC,
		//},
		{
			name:   "rest",
			server: ServeHttp,
		},
	}
)

func StartServers(cfgs *viper.Viper, servers []ServerInfo) error {
	for _, server := range servers {
		cancelFunChan := make(chan CancelFun, 1)
		go server.server(cancelFunChan, cfgs)
		server.cancel = <-cancelFunChan
		log.Printf("cancel func... %v", server.cancel)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c

	log.Printf("Stopping service with sig %v", sig)

	for i := len(servers) - 1; i >= 0; i-- {
		server := servers[i]
		log.Println(server.name)
		log.Println(server.cancel)
		if err := server.cancel(); err != nil {
			log.Printf("Stop server %v error %v", server.name, err)
			return err
		}
	}
	return nil
}
