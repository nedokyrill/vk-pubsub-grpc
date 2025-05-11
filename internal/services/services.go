package services

import "github.com/gin-gonic/gin"

type PubSubService interface {
	SubscribeFromRPC(c *gin.Context)
	PublishFromRPC(c *gin.Context)
}
