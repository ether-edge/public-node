package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	
	"public-node/internal/config"

)

func SendPostRequest() {
	baseUrl := config.GetEnvData("base_url")
	if baseUrl == "" {
		baseUrl = "http://127.0.0.1:8080" // Set default value if not found
	}

	apiURL := fmt.Sprintf("%s/api/public-node-activity/store", baseUrl)

	formValues := url.Values{}
	formValues.Set("LicenseID",config.GetEnvData("license_id"))
	formValues.Set("ActivityDate", time.Now().Format("2006-01-02"))
	formValues.Set("HoursOnline", strconv.Itoa(1))

	payload := bytes.NewBufferString(formValues.Encode())

	req, err := http.NewRequest("POST", apiURL, payload)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("License", config.GetEnvData("license"))
	req.Header.Set("Api-Key", config.GetEnvData("api_key"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("POST request failed with status: %d, response: %s\n", resp.StatusCode, body)
		return
	}

	// fmt.Println("POST request successful:", string(body))
}