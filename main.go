package main

import (
	"fmt"
	"time"

	"fileUpload/internal/config"
	"fileUpload/internal/handlers"

	"github.com/gofiber/fiber/v2"
)


func main() {

	config.GetInputForAPISection()

	app := fiber.New(fiber.Config{
		Prefork:      true, // Enable multiple OS processes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    1 * 1024 * 1024 * 1024, // 1 GB
	})
	app.Post("/upload-file", handlers.CreateFileHandler)
	app.Get("/test", handlers.TestHandler)
	app.Get("/download/:filename", handlers.DownloadFileHandler)

	fmt.Println("Server is running on port 40252...")
	app.Listen(":40252")
}
