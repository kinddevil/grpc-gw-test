package servers

import (
	"github.com/spf13/viper"
	"log"
)

func serverProm(terminate chan<- CancelFun, cfgs *viper.Viper) {
	promPort := cfgs.GetString("monitor.prom.port")
	log.Println(promPort)
}
