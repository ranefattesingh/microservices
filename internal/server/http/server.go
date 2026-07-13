package http

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ranefattesingh/ecommerce-platform/internal/config"
)

type httpServer struct {
	config         config.ServerConfig
	shutdownSignal chan struct{}
}

func NewHTTPServer(conf config.ServerConfig) *httpServer {
	return &httpServer{
		config:         conf,
		shutdownSignal: make(chan struct{}, 1),
	}
}

func (s *httpServer) StartServer() error {
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(s.config.Port),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	go s.shutdown(srv)

	err := srv.ListenAndServe()

	<-s.shutdownSignal

	if err != nil && err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *httpServer) shutdown(srv *http.Server) {
	defer close(s.shutdownSignal)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		slog.Error("server shutdown", slog.Any("error", err))
	}

	slog.Info("server exiting")
}
