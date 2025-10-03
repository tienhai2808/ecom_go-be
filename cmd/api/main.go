package main

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/server"

	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Khởi tạo server thất bại: %v", err)
	}

	ch := make(chan error, 1)

	go func() {
		if err := server.Start(); err != nil {
			ch <- err
		}
	}()

	log.Println("Chạy server thành công")

	server.GracefulShutdown(ch)
}
