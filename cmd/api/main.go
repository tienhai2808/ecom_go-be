package main

import "backend/internal/server"

func main() {
	application := server.NewApplication()
	defer application.Close()

	application.Run()
}
