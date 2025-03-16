package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func SendPostRequest() (string, error) {

	apiURL := "http://127.0.0.1:8080/api/public-node-activity/store"

	formValues := url.Values{}
	formValues.Set("LicenseID", "a4e35532-9c6a-4b65-8ec0-0c9a39436e0b")
	formValues.Set("ActivityDate", time.Now().Format("2006-01-02"))
	formValues.Set("HoursOnline", strconv.Itoa(1))

	payload := strings.NewReader(formValues.Encode())

	req, err := http.NewRequest("POST", apiURL, payload)
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("POST request failed with status: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error decoding API response: %v", err)
	}

	return string(body), nil
}

