package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/ranefattesingh/microservices/user/config"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/microservices/user/repo/psql"
	"github.com/ranefattesingh/microservices/user/route"

	"github.com/ranefattesingh/pkg/configloader"
	"github.com/ranefattesingh/pkg/log"
	pgmigrate "github.com/ranefattesingh/pkg/postgresql/migrate"
	"github.com/ranefattesingh/pkg/postgresql/pgx/pool"
	httpsrv "github.com/ranefattesingh/pkg/server/http"
)

const (
	migrationDir = "repo/psql/migrations"
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
		Encoder:  log.DefaultJSONEncoder(),
	})

	dbPool, err := pool.NewDatabaseConnectionPool(ctx, configuration.DatabaseConfig.GetConnectionStringWithOptions())
	if err != nil {
		return fmt.Errorf("database connectivity error: %w", err)
	}

	log.Logger().Info("database connection established")

	if configuration.GinMode == gin.DebugMode {
		migrator, err := pgmigrate.NewDatabaseMigrator(migrationDir, configuration.DatabaseConfig.GetConnectionString())
		if err != nil {
			return fmt.Errorf("migrator initialization error: %w", err)
		}

		err = migrator.Migrate().Up()
		if !errors.Is(err, migrate.ErrNoChange) && err != nil {
			return fmt.Errorf("migration up error: %w", err)
		}

		log.Logger().Info("database migration complete")
	}

	repo := psql.NewRepo(dbPool)
	service := core.NewUserService(repo)

	router := route.NewRouter(configuration.GinMode)

	router.Handle(service)

	srv := httpsrv.NewHTTPServer(
		configuration.HTTPConfig.Host,
		configuration.HTTPConfig.Port,
		func() error {
			dbPool.CloseDatabaseConnectionPool()

			return nil
		},
	)

	log.Logger().Info("server starting on port", zap.Int("port", configuration.HTTPConfig.Port))

	if err := srv.Start(router, cancelFn); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
