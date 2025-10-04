package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saikewei/my_website/back/internal/config"
	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/photo"
)

func main() {
	config.LoadConfig()

	database.Connect()

	r := gin.Default()

	apiGroup := r.Group("/api")

	photo.RegisterRouters(apiGroup)

	r.Run(":9000") // 监听
}
