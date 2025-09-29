package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saikewei/my_website/back/photo"
)

func main() {
	r := gin.Default()

	photo.RegisterRouters(r)

	r.Run(":9000") // 监听并在
}
