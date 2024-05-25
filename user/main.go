package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ranefattesingh/microservices/user/config"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/microservices/user/route"
	"github.com/ranefattesingh/pkg/configloader"
	"github.com/ranefattesingh/pkg/log"
	"github.com/ranefattesingh/pkg/pgx/pool"
	httpSrv "github.com/ranefattesingh/pkg/server/http"
	"go.uber.org/zap"
)

func main() {
	err := run(context.WithCancel(context.Background()))
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	log.Logger().Info("server shutdown complete")
}

func run(ctx context.Context, cancelFn context.CancelFunc) error {
	var configuration config.MainConfig

	err := configloader.NewDefaultLoader().Load(ctx, &configuration)
	if err != nil {
		return fmt.Errorf("configloader error: %w", err)
	}

	log.Init(log.Config{
		Output:   os.Stdout,
		LogLevel: log.LogLevel(configuration.LogLevel),
	})

	pool, err := pool.NewDatabaseConnectionPool(ctx, configuration.DatabaseConfig.GetConnectionString())
	if err != nil {
		return fmt.Errorf("database connectivity error: %w", err)
	}

	log.Logger().Info("database connection established")

	service := core.NewUserService()

	r := route.NewRouter(configuration.GinMode)

	r.Handle(service)

	srv := httpSrv.NewHTTPServer(
		configuration.HTTPConfig.Host,
		configuration.HTTPConfig.Port,
		func() error {
			pool.CloseDatabaseConnectionPool()
			return nil
		},
	)

	log.Logger().Info("server starting on port", zap.Int("port", configuration.HTTPConfig.Port))
	if err := srv.Start(r, cancelFn); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
