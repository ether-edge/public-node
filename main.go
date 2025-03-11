package main

import (

	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)


func createFileHandler(c *fiber.Ctx) error {
    fileHeader, err := c.FormFile("uploaded_file")
    if err != nil {
        fmt.Println("Error getting file:", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
    }

    fileName := c.FormValue("filename", fileHeader.Filename)
    
    folderName := "uploads"
    if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
        fmt.Println("Error creating folder:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create folder"})
    }

    filePath := fmt.Sprintf("%s/%s", folderName, fileName)
    if err := c.SaveFile(fileHeader, filePath); err != nil {
        fmt.Println("Error saving file:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
    }

    fmt.Println("File uploaded successfully:", filePath)
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message":  "File uploaded successfully!",
        "filename": fileName,
        "path":     filePath,
    })
}


func testHandler(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"message": "File uploaded successfully!"})
}

func main() {

	app := fiber.New(fiber.Config{
		Prefork:      true, // Enable multiple OS processes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		BodyLimit:    1 * 1024 * 1024 * 1024, // 1 GB
	})
	app.Post("/upload-file", createFileHandler)
	app.Get("/test", testHandler)

	fmt.Println("Server is running on port 40252...")
	app.Listen(":40252")
}
