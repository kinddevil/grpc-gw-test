package servers

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var (
	promRegistry = prometheus.NewRegistry()
)

func ServerProm(terminate chan<- CancelFun, cfgs *viper.Viper) {
	promPort := cfgs.GetString("common.prom.port")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register prometheus and open http port for scraping
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
		Addr:    promPort,
	}

	terminate <- func() error {
		return httpServer.Shutdown(ctx)
	}

	log.Printf("start prometheus service with %v", promPort)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server for prometheus.")
		}
	}()
}
