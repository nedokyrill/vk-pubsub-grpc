package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nedokyrill/vk-pubsub-grpc/internal/api/pubSubHandlers"
	httpserver "github.com/nedokyrill/vk-pubsub-grpc/internal/httpServer"
	"github.com/nedokyrill/vk-pubsub-grpc/internal/services/pubSubService"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	pb "github.com/nedokyrill/vk-pubsub-grpc/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunApp() {
	// Init logger
	logger.InitLogger()
	defer logger.Logger.Sync()

	// Load environment
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading .env file")
	}

	// Connection with remote RPC server
	// Instead "localhost" you can use name of the docker container, which contains gRPC service
	conngRPC, err := grpc.NewClient("pub-sub:"+os.Getenv("GRPC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Error("did not connect to gRPC server: %v", err)
	}
	defer conngRPC.Close()

	clientRPC := pb.NewPubSubClient(conngRPC)

	// Init services
	pubSubServ := pubSubService.NewPubSubService(clientRPC)

	// Init handlers
	pubSubHand := pubSubHandlers.NewPubSubHandler(pubSubServ)

	// Init default route
	r := gin.Default()
	api := r.Group("/api/v1")

	// Register routes
	pubSubHand.RegisterRoutes(api)

	// Init n conf API server
	srv := httpserver.NewAPIServer(r)

	// Start server
	go srv.Start()
	logger.Logger.Info("API server started successfully on port " + os.Getenv("API_PORT"))

	// Shutdown server (with graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalw("Shutdown error",
			"error", err)
	}
}
