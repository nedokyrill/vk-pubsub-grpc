package pubSubService

import (
	"github.com/gin-gonic/gin"
	"github.com/nedokyrill/vk-pubsub-grpc/pkg/logger"
	"net/http"
)

type SubscribeRequest struct {
	Key string `json:"key"`
}

func (s *PubSubService) SubscribeFromRPC(c *gin.Context) {
	var payload SubscribeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Logger.Errorf("Error with unmarshalling JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
