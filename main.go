package main

import (
	"fmt"
	"time"
    "os"

	"fileUpload/internal/config"
	"fileUpload/internal/handlers"

	"github.com/gofiber/fiber/v2"
)


func main() {
    port := "40252"

	// if !config.IsPortAvailable(port) {
	// 	fmt.Println("Port", port, "is occupied. Exiting.")
	// 	os.Exit(1)
	// }

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
	err := app.Listen(":" + port)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
