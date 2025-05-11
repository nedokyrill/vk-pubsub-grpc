package pubSubService

import pb "github.com/nedokyrill/vk-pubsub-grpc/pkg/proto"

type PubSubService struct {
	PubSubServiceClient pb.PubSubClient
}

func NewPubSubService(pubSubServiceClient pb.PubSubClient) *PubSubService {
	return &PubSubService{
		PubSubServiceClient: pubSubServiceClient,
	}
}
