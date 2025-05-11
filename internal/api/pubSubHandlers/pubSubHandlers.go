package pubSubHandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nedokyrill/vk-pubsub-grpc/internal/services"
)

type PubSubHandler struct {
	PubSubService services.PubSubService
}

func NewPubSubHandler(service services.PubSubService) *PubSubHandler {
	return &PubSubHandler{
		PubSubService: service,
	}
}

func (h *PubSubHandler) RegisterRoutes(router *gin.RouterGroup) {
	pubSubRouter := router.Group("/pub-sub")
	{
		pubSubRouter.POST("/subscribe", h.PubSubService.SubscribeFromRPC)
		pubSubRouter.POST("/publish", h.PubSubService.PublishFromRPC)
	}
}
