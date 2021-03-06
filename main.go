package main

import (
	"fmt"

	"net/http"

	"github.com/frankffenn/aquarium/config"
	"github.com/frankffenn/aquarium/routers"
	"github.com/frankffenn/aquarium/utils/env"
	"github.com/frankffenn/aquarium/utils/log"
	glog "github.com/frankffenn/aquarium/utils/log"
)

func main() {
	configJSON := "config.json"
	if err := config.InitConfig(configJSON); err != nil {
		log.Fatal("init config failed, %v", err)
	}

	logLevel := "info"
	if config.Configs.RunMode == "debug" {
		logLevel = "debug"
	}

	logger, _ := glog.NewLogger(env.LogPath, logLevel)
	glog.SetDefault(logger)

	router := routers.InitRouter()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Configs.Host, config.Configs.Port),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("listen: %s\n", err)
	}
}
