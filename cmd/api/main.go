package main

import (
	"backend/config"
	"backend/internal/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		if err := server.Start(); err != nil {
			errCh <- err
		}
	}()

	log.Println("Chạy server thành công")

	select {
	case err = <-errCh:
		log.Printf("Chạy server thất bại: %v", err)
	case <-stop:
		log.Println("Có tín hiệu dừng server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
