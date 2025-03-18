package main

import (
	"fmt"
	"os"
	"time"

	"public-node/internal/config"
	"public-node/internal/handlers"
	"public-node/internal/jobs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := "40252"

	config.GetInputForAPISection()

	success, err := config.MakeAPICall()
	if err != nil {
		fmt.Println("Error during API call:", err) // Todo :: change This massage
		os.Exit(1)
	}
	if !success {
		fmt.Println("API call unsuccessful.")
		os.Exit(1) // Exit if the API call was unsuccessful
	}

	//Todo :: Fix single core use for corn only
	jobs.RunCronJobsStarted()

	app := fiber.New(fiber.Config{
		Prefork:      false,   
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    1 * 1024 * 1024 * 1024, // 1 GB
	})
	app.Post("/upload-file", handlers.CreateFileHandler)
	app.Get("/test", handlers.TestHandler)
	app.Get("/download/:filename", handlers.DownloadFileHandler)

	fmt.Println("Server is running on port 40252...")
	err = app.Listen(":" + port)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}

}
