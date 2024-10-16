package main

import (
	server "API_Call_GVA_and_LLMAI/cmd"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Kiểm tra và nạp biến môi trường
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // Cổng mặc định
	}

	// Khởi tạo ApiServer
	apiServer := server.NewApiServer()

	// In ra cổng đang sử dụng
	log.Printf("Starting server on port %s...", port)

	// Chạy server trong goroutine để có thể tắt khi nhận tín hiệu ngắt
	go apiServer.Start(port) // Chỉ gọi hàm mà không cần xử lý giá trị trả về

	// Chờ tín hiệu ngắt (Ctrl+C) hoặc tín hiệu từ hệ thống (kill signal)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Dừng server một cách nhẹ nhàng
	apiServer.Stop()

	log.Println("Server stopped.")
}
