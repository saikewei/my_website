package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/saikewei/my_website/back/auth"
	"github.com/saikewei/my_website/back/internal/config"
	"github.com/saikewei/my_website/back/internal/database"
	"github.com/saikewei/my_website/back/photo"
)

func main() {
	config.LoadConfig()

	database.Connect()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // 允许你的前端源
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	apiGroup := r.Group("/api")

	photo.RegisterRouters(apiGroup)
	auth.RegisterAuthRouters(apiGroup)

	r.Run(":9000") // 监听
}
