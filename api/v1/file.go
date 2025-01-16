package v1

import (
	"fitbyte/internal/handler"
	"fitbyte/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {

	protected := router.Group("file")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/file", handler.UploadFileHandler)
		protected.GET("/file/:id", handler.GetFileHandler)
	}

}
