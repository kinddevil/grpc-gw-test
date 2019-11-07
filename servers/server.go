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
	Name   string
	Server Server
	Cancel CancelFun
}

var (
	ServerList = []*ServerInfo{
		{
			Name:   "grpc",
			Server: ServeGRPC,
		},
		{
			Name:   "rest",
			Server: ServeHttp,
		},
		{
			Name:   "prometheus",
			Server: ServeProm,
		},
	}
)

func StartServers(cfgs *viper.Viper, servers []*ServerInfo) error {
	for _, server := range servers {
		cancelFunChan := make(chan CancelFun, 1)
		go server.Server(cancelFunChan, cfgs)
		server.Cancel = <-cancelFunChan
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGCHLD)
	sig := <-c

	log.Printf("Stopping service with sigs %v", sig)

	for i := len(servers) - 1; i >= 0; i-- {
		server := servers[i]
		if err := server.Cancel(); err != nil {
			log.Printf("Stop server %v error %v", server.Name, err)
			return err
		}
	}
	return nil
}
