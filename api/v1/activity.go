package v1

import (
	"fitbyte/internal/handler"
	"fitbyte/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(r *gin.RouterGroup) {

	protected := r.Group("/activity")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/", handler.CreateActivity)
	}
}

//Kegunaan Protected itu ntar buat kalau mau akses itu harus login
