package v1

import (
	"fitbyte/internal/handler"
	"fitbyte/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {

	{
		router.POST("/login", handler.LoginUser)
		router.POST("/register", handler.RegisterUser)
	}
	//Kegunaan Protected itu ntar buat kalau mau akses itu harus login
	protected := router.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		/*
				USE EXAMPLE
			protected.GET("/", handler.GetUserProfileHandler)
		*/
	}
}
