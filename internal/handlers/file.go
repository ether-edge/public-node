package handlers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
)


func CreateFileHandler(c *fiber.Ctx) error {
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


func DownloadFileHandler(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Could not load .env file, using default values")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	fileName := c.Params("filename")
	filePath := fmt.Sprintf("uploads/%s", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Error: File not found:", filePath)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	downloadURL := fmt.Sprintf("%s/download/%s", baseURL, fileName)

	return c.SendFile(downloadURL)
}

func TestHandler(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"message": "File uploaded successfully!"})
}