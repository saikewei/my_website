package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 定义一个简单的 API 路由
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头，允许跨域请求 (开发时非常重要)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "Hello from Go Backend!")
	})

	fmt.Println("Backend server is running on http://localhost:8080")
	// 启动服务器，监听 8080 端口
	log.Fatal(http.ListenAndServe(":8080", nil))
}
