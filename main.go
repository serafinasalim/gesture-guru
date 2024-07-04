package main

import (
	"github.com/gin-gonic/gin"
	"github.com/serafinasalim/gesture-guru/controllers/userController"
	"github.com/serafinasalim/gesture-guru/models"
)

func main() {
	route := gin.Default()
	models.ConnectDatabase()

	route.GET("gesture-guru/users/:id", userController.Detail)
	route.POST("gesture-guru/users/register", userController.Register)
	route.POST("gesture-guru/users/login", userController.Login)
	// route.PUT("gesture-guru/users/")

	route.Run()
}
