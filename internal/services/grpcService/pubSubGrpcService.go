package grpcService

import (
	"context"
	"fmt"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	pb "github.com/nedokyrill/vk-pubsub-grpc/pkg/proto"
	"github.com/nedokyrill/vk-pubsub-grpc/subpub"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	pb.UnimplementedPubSubServer
	pubSub subpub.SubPub
}

func (s *GrpcServer) Subscribe(in *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.Event]) error {
	subject := in.GetKey()
	ctx := stream.Context()

	logger.Logger.Info("Subscribe on subject ", subject, " started successfully")

	sub, err := s.pubSub.Subscribe(subject, func(msg interface{}) {
		event, ok := msg.(*pb.Event)
		if !ok {
			logger.Logger.Error("Invalid event type")
			return
		}
		if err := stream.Send(event); err != nil {
			logger.Logger.Error("Stream send failed on subject ", subject, ", error: ", err)
		}
	})
	if err != nil {
		logger.Logger.Error("Subscription failed on subject ", subject, ", error: ", err)
		return err
	}

	<-ctx.Done()
	logger.Logger.Info("Stream context done on subject ", subject)
	sub.(*subpub.SubscriptionAgent).Done <- struct{}{}

	return ctx.Err()
}

func (s *GrpcServer) Publish(_ context.Context, in *pb.PublishRequest) (*emptypb.Empty, error) {
	subject := in.GetKey()
	data := in.GetData()
	event := &pb.Event{Data: data}

	err := s.pubSub.Publish(subject, event)
	if err != nil {
		logger.Logger.Errorf("Failed to publish event: %v", err)
		return nil, fmt.Errorf("publish failed: %w", err)
	}

	logger.Logger.Info("Successfully published event " + in.Data + " in subject " + in.Key)
	return &emptypb.Empty{}, nil
}
