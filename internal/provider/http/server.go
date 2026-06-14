package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"identify/internal/config"
	"identify/internal/handler"
	"identify/internal/logger"
)

type HTTPProvider struct {
	server   *http.Server
	listener net.Listener
	config   config.AppConfig
	logger   logger.Logger
}

// NewHTTPProvider creates a new HTTPProvider instance.
func NewHTTPProvider(cfg config.AppConfig, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, logger logger.Logger) *HTTPProvider {
	engine := gin.New()
	SetupRoutes(engine, userHandler, authHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	return &HTTPProvider{
		server: server,
		config: cfg,
		logger: logger,
	}
}

func (h *HTTPProvider) Name() string {
	return "http"
}

// Boot starts HTTP server. It binds to the TCP port synchronously to catch startup errors,
// then serves requests in a separate goroutine.
func (h *HTTPProvider) Boot(ctx context.Context) error {
	h.logger.Info("starting HTTP server", "port", h.config.Port)

	listener, err := net.Listen("tcp", h.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind port %d: %w", h.config.Port, err)
	}
	h.listener = listener

	go func() {
		if err := h.server.Serve(h.listener); err != nil && err != http.ErrServerClosed {
			h.logger.Error("HTTP server error", "error", err)
		}
	}()

	return nil
}

// Close gracefully closes the HTTP server.
func (h *HTTPProvider) Close(ctx context.Context) error {
	h.logger.Info("shutting down HTTP server")
	return h.server.Shutdown(ctx)
}
