package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
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
func MakeAPICall() (bool, error) {

	url := "http://127.0.0.1:8080/node-status" // API URL

	license := GetEnvData("license") 
	if license == "" {
		return false, fmt.Errorf("license not found in environment variables")
	}

	apiValuesWithLicense := map[string]string{
		"license": license,
	}

	reqBody, err := json.Marshal(apiValuesWithLicense)
	if err != nil {
		return false, fmt.Errorf("error encoding JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API call failed with status: %d", resp.StatusCode) // Todo :: change This massage 
	}

	var response map[string]interface{} 
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, fmt.Errorf("error decoding API response: %v", err)
	}

	success, exists := response["status"]
	if !exists {
		return false, fmt.Errorf("invalid API response format: no 'status' field found")
	}

	if status, ok := success.(bool); ok {
		return status, nil
	} else {
		return false, fmt.Errorf("'status' field is not a boolean in the response")
	}
}

func LoadEnv() error {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}
	return nil
}

// GetEnvData retrieves the value of an environment variable
func GetEnvData(varName string) string {
	// Load environment variables from the .env file if not already loaded
	if err := LoadEnv(); err != nil {
		log.Fatal(err)
	}

	// Retrieve the environment variable by name
	value := os.Getenv(varName)

	if value == "" {
		fmt.Printf("Environment variable %s is not set\n", varName)
	}

	return value
}

func IsPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Port", port, "is already in use.")
		return false
	}
	defer listener.Close()
	return true
}
