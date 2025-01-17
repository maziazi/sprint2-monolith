package v1

import (
	"fitbyte/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(router *gin.RouterGroup) {
	router.POST("/activity", handler.CreateActivity)
}
