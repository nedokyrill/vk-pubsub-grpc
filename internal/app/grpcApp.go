package app

import (
	"github.com/joho/godotenv"
	"github.com/nedokyrill/vk-pubsub-grpc/internal/services/grpcService"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	pb "github.com/nedokyrill/vk-pubsub-grpc/pkg/proto"
	"google.golang.org/grpc"
	"net"
	"os"
)

func GrpcRun() {
	// Инициализация логгера
	logger.InitLogger()
	defer logger.Logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading .env file")
	}
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		logger.Logger.Fatal("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, &grpcService.GrpcServer{})
	logger.Logger.Info("Listening gRPC server at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Logger.Fatal("Failed to serve:", err)
	}

}
