package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/ranefattesingh/microservices/user/config"
	"github.com/ranefattesingh/microservices/user/route"
	"github.com/ranefattesingh/pkg/configloader"
	"github.com/ranefattesingh/pkg/log"
	httpSrv "github.com/ranefattesingh/pkg/server/http"
	"go.uber.org/zap"
)

func main() {
	appCtx, cancelFn := context.WithCancel(context.Background())

	var configuration config.MainConfig

	err := configloader.NewDefaultLoader().Load(appCtx, &configuration)
	if err != nil {
		panic("config loading failed: " + err.Error())
	}

	log.Init(log.Config{
		Output:   os.Stdout,
		LogLevel: log.LogLevel(configuration.LogLevel),
	})

	r := route.NewRouter(configuration.GinMode)

	r.Handle()

	srv := httpSrv.NewHTTPServer(configuration.HTTPConfig.Host, configuration.HTTPConfig.Port)

	shutdownCh := make(chan struct{}, 1)
	go srv.Shutdown(shutdownCh, cancelFn)

	err = srv.Start(r)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Error("server exited with error", zap.Error(err))
	}

	<-shutdownCh
}
