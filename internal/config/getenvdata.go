package config

import (
	"bufio"
	"fmt"
	"os"
	
	"strings"
)

func GetInputForAPISection() map[string]string {
	apiValues := make(map[string]string)

	existingValues := readEnvFile()

	if existingValues["license"] == "" {
		fmt.Println("Enter your Node credentials")

		fmt.Print("License: ")
		apiValues["license"] = promptInput("license", existingValues)

		fmt.Print("API Key: ")
		apiValues["api_key"] = promptInput("api_key", existingValues)

		fmt.Print("Address: ")
		apiValues["address"] = promptInput("address", existingValues)

		fmt.Print("Base URL: ")
		apiValues["base_URL"] = promptInput("base_URL", existingValues)

		writeToEnvFile(apiValues)
	} else {
		apiValues = existingValues
	}

	return apiValues
}

func readEnvFile() map[string]string {
	existingValues := make(map[string]string)

	if _, err := os.Stat(".env"); err == nil {
		file, err := os.Open(".env")
		if err != nil {
			fmt.Println("Error opening .env file:", err)
			return existingValues
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existingValues[parts[0]] = parts[1]
			}
		}
	}

	return existingValues
}

func promptInput(key string, existingValues map[string]string) string {
	if value, exists := existingValues[key]; exists && value != "" {
		return value
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		return scanner.Text()
	}
}

func writeToEnvFile(apiValues map[string]string) {
	existingValues := readEnvFile()

	for key, value := range apiValues {
		existingValues[key] = value
	}

	file, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening .env file:", err)
		return
	}
	defer file.Close()

	for key, value := range existingValues {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			fmt.Println("Error writing to .env file:", err)
			return
		}
	}
}
