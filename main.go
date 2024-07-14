package main

import (
	"github.com/gin-gonic/gin"
	lessonController "github.com/serafinasalim/gesture-guru/controllers/lessonController"
	userController "github.com/serafinasalim/gesture-guru/controllers/userController"
	"github.com/serafinasalim/gesture-guru/models"

	_ "github.com/serafinasalim/gesture-guru/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gesture Guru API
// @version 1.0
// @description Documentation for  GestureGuru

// @host localhost:8080
// @BasePath /gesture-guru

func main() {
	route := gin.Default()
	models.ConnectDatabase()

	// Swagger setup
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := route.Group("/gesture-guru")
	{
		// User API
		api.GET("/user/:id", userController.Detail)
		api.POST("/user/register", userController.Register)
		api.PUT("/user/verify/:id", userController.Verify)
		api.PUT("/user/request-otp/:id", userController.RequestOTP)
		api.POST("/user/login", userController.Login)
		api.PUT("/user/:id", userController.UpdateUser)
		api.PUT("/user/email/:id", userController.UpdateEmail)

		// Lesson API
		api.GET("/lessons", lessonController.Browse)
	}

	route.Run()
}
