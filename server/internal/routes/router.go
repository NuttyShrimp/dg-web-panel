package routes

import (
	"degrens/panel/lib/log"

	"github.com/gin-gonic/gin"
)

type Router struct {
	RouterGroup *gin.RouterGroup
	Logger      log.Logger
}
