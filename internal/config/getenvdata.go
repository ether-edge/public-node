package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

)

func GetInput() map[string]string {
	apiValues := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter your Node credentials")

	fmt.Print("Address: ")
	scanner.Scan()
	apiValues["address"] = scanner.Text()

	fmt.Print("License: ")
	scanner.Scan()
	apiValues["license"] = scanner.Text()

	fmt.Print("API key: ")
	scanner.Scan()
	apiValues["api_key"] = scanner.Text()


	fmt.Print("NODE URL: ")
	scanner.Scan()
	apiValues["node_url"] = SanitizeURL(scanner.Text())


	fmt.Print("BASE URL: ")
	scanner.Scan()
	apiValues["base_url"] = SanitizeURL(scanner.Text())

	setLinuxEnvironmentVariables(apiValues)

	return apiValues
}

func setLinuxEnvironmentVariables(apiValues map[string]string) {
	unsetEnvIfExists("ADDRESS")
	unsetEnvIfExists("LICENSE")
	unsetEnvIfExists("API_KEY")
	unsetEnvIfExists("NODE_URL")
	unsetEnvIfExists("BASE_URL")
	
	os.Setenv("ADDRESS", apiValues["address"])
	os.Setenv("LICENSE", apiValues["license"])
	os.Setenv("API_KEY", apiValues["api_key"])
	os.Setenv("NODE_URL", apiValues["node_url"])
	os.Setenv("BASE_URL", apiValues["base_url"])

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Error getting home directory:", err)
		return
	}

	bashrcFile := homeDir + "/.bashrc"

	file, err := os.OpenFile(bashrcFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("❌ Error opening .bashrc:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "\n# Added by Go program\nexport ADDRESS=%s\nexport LICENSE=%s\nexport API_KEY=%s\nexport NODE_URL=%s\nexport BASE_URL=%s\n",
	apiValues["address"], apiValues["license"], apiValues["api_key"], apiValues["node_url"], apiValues["base_url"])
if err != nil {
	fmt.Println("❌ Error appending to .bashrc:", err)
	return
}

}


func unsetEnvIfExists(key string) {
	if _, exists := os.LookupEnv(key); exists {
		os.Unsetenv(key)
		fmt.Printf("❌ Unset %s from current session\n", key)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Error getting home directory:", err)
		return
	}
	bashrcFile := homeDir + "/.bashrc"

	content, err := ioutil.ReadFile(bashrcFile)
	if err != nil {
		fmt.Println("❌ Error reading .bashrc:", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "export "+key+"=") {
			continue
		}
		newLines = append(newLines, line)
	}

	updatedContent := strings.Join(newLines, "\n")
	err = ioutil.WriteFile(bashrcFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Println("❌ Error updating .bashrc:", err)
		return
	}
}


func GetEnvData(varName string) string {
	
	envVars := GetBashrcEnv()

	value := "not set"

	if val, exists := envVars[varName]; exists {
		value = val
	}

	return value
}

func GetBashrcEnv() map[string]string {
	envData := make(map[string]string)
	file, err := os.Open(os.Getenv("HOME") + "/.bashrc")
	if err != nil {
		fmt.Println("Error opening .bashrc:", err)
		return envData
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "export ") {
			parts := strings.SplitN(strings.TrimPrefix(line, "export "), "=", 2)
			if len(parts) == 2 {
				envData[parts[0]] = strings.Trim(parts[1], `"`)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .bashrc:", err)
	}

	return envData
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

func WriteSingleToBashrc(key string, value string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Error getting home directory:", err)
		return
	}

	bashrcFile := homeDir + "/.bashrc"

	content, err := ioutil.ReadFile(bashrcFile)
	if err != nil {
		fmt.Println("❌ Error reading .bashrc:", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	keyPrefix := "export " + key + "="

	for _, line := range lines {
		if !strings.HasPrefix(line, keyPrefix) {
			newLines = append(newLines, line)
		}
	}

	newLines = append(newLines, keyPrefix+value)

	updatedContent := strings.Join(newLines, "\n")

	err = ioutil.WriteFile(bashrcFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Println("❌ Error writing to .bashrc:", err)
		return
	}
}

func SanitizeURL(url string) string {
	url = strings.TrimSpace(url) 
	url = strings.TrimSuffix(url, "/")
	url = strings.ReplaceAll(url, " ", "")
	return url
}

func MakeAPICall() (bool, error) {
	baseUrl := GetEnvData("BASE_URL")

	url := fmt.Sprintf("%s/node-status", baseUrl)

	license := GetEnvData("license")
	if license == "" {
		return false, fmt.Errorf("license not found in environment variables")
	}

	nodeUrl := GetEnvData("node_url")
	if license == "" {
		return false, fmt.Errorf("Node Url not found in environment variables")
	}

	apiValuesWithLicense := map[string]string{
		"license":  license,
		"node_url": nodeUrl,
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
		if licenseID, ok := response["license"].(string); ok {
			WriteSingleToBashrc("LICENSE_ID", licenseID)
		} else {
			return false, fmt.Errorf("'license' field is not a string in the response")
		}
		return status, nil
	} else {
		return false, fmt.Errorf("'status' field is not a boolean in the response")
	}
}