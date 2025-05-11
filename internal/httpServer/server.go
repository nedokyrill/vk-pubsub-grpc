package httpserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	"net/http"
	"os"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(router *gin.Engine) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         ":" + os.Getenv("API_PORT"),
			Handler:      router.Handler(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *APIServer) Start() {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatal("Error with starting API server: ", err)
	}
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	switch ctx.Err() {
	case context.DeadlineExceeded:
		logger.Logger.Error("Timeout shutting down server")
	case nil:
		logger.Logger.Info("Shutdown completed before timeout.")
	default:
		logger.Logger.Error("Shutdown ended with:", ctx.Err())
	}

	return nil
}
