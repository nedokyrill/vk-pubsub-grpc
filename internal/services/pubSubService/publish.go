package pubSubService

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	pb "github.com/nedokyrill/vk-pubsub-grpc/pkg/proto"
	"net/http"
	"time"
)

type PublishPayload struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

func (s *PubSubService) PublishFromRPC(c *gin.Context) {
	var payload PublishPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Logger.Errorf("Error with unmarshalling JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	req := &pb.PublishRequest{
		Key:  payload.Key,
		Data: payload.Data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := s.PubSubServiceClient.Publish(ctx, req)
	if err != nil {
		logger.Logger.Errorf("Error publishing message: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logger.Logger.Infof("Published event via gRPC: key=%s, data=%s", payload.Key, payload.Data)
	c.JSON(http.StatusOK, gin.H{})

}
