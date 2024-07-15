package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	lessonController "github.com/serafinasalim/gesture-guru/controllers/lessonController"
	userController "github.com/serafinasalim/gesture-guru/controllers/userController"
	"github.com/serafinasalim/gesture-guru/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gesture Guru API
// @version 1.0
// @description Documentation for GestureGuru

// @host localhost:8080
// @BasePath /gesture-guru

func main() {
	route := gin.Default()
	// export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

	// Middleware CORS
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://a08f-114-122-7-105.ngrok-free.app/"}, // Ubah sesuai kebutuhan
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	models.ConnectDatabase()

	// Swagger setup
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := route.Group("/gesture-guru")
	{
		// User API
		api.POST("/user", userController.Detail)
		api.POST("/user/register", userController.Register)
		api.PUT("/user/verify/:id", userController.Verify)
		api.PUT("/user/request-otp/:id", userController.RequestOTP)
		api.POST("/user/login", userController.Login)
		api.PUT("/user/:id", userController.UpdateUser)
		api.PUT("/user/email/:id", userController.UpdateEmail)

		// Lesson API
		api.POST("/lessons", lessonController.Browse)
		api.POST("/lesson", lessonController.Detail)
		api.PUT("/lesson/save/:userId/:lessonId", lessonController.SaveLesson)
		api.PUT("/lesson/attempt/:userId/:lessonId", lessonController.AttemptLesson)
	}

	route.Run()
}
