package main

import (
	"github.com/gin-gonic/gin"
	"github.com/serafinasalim/gesture-guru/controllers/userController"
	"github.com/serafinasalim/gesture-guru/models"
)

func main() {
	route := gin.Default()
	models.ConnectDatabase()

	route.GET("gesture-guru/user/:id", userController.Detail)

	route.Run()
}
