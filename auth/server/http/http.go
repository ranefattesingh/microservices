package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/ranefattesingh/microservices/auth/config"
	"github.com/ranefattesingh/microservices/auth/router"
)

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
)

type httpServer struct {
	logger *zap.Logger
	server *http.Server

	addr              string
	readTimeout       time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
}

func NewHTTPServer(
	logger *zap.Logger,
	c config.ServerConfig,
) *httpServer {
	return &httpServer{
		logger:            logger.Named("http-server"),
		addr:              fmt.Sprintf("%s:%d", c.Host, c.Port),
		readTimeout:       defaultReadTimeout,
		idleTimeout:       defaultIdleTimeout,
		readHeaderTimeout: defaultReadHeaderTimeout,
		writeTimeout:      defaultWriteTimeout,
	}
}

func (h *httpServer) Start(providers ...router.RouteProvider) error {
	r := router.NewRouter(providers...)

	h.server = &http.Server{
		Addr:              h.addr,
		ReadHeaderTimeout: h.readHeaderTimeout,
		ReadTimeout:       h.readTimeout,
		WriteTimeout:      h.writeTimeout,
		IdleTimeout:       h.idleTimeout,
		Handler:           r,
	}

	h.logger.Info("starting HTTP server", zap.String("addr", h.addr))

	if err := h.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("HTTP server: %w", err)
	}

	return nil
}

func (h *httpServer) Shutdown() error {
	if h.server == nil {
		return nil
	}

	h.logger.Info("shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Error("failed to shutdown HTTP server", zap.Error(err))

		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	h.logger.Info("HTTP server stopped")

	return nil
}
