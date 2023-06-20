package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	apiURL     = "http://13.235.66.99:8089/agent-instructions?hostIP=%s"
	orgID      = "PRA646ef244a79a86633dffebd4"
	hostIP     = "192.168.6.40"
	pollPeriod = 10 * time.Second // Polling interval in seconds
)

func main() {
	lastUpdateTime := time.Time{}

	for {
		// Prepare the request
		url := fmt.Sprintf(apiURL, hostIP)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			continue
		}

		// Set headers
		req.Header.Set("company-code", orgID)

		// Send the request
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request failed:", err)
			continue
		}

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("cannot read response body, error: %+v", err)
			return
		}

		ioutil.WriteFile("response.json", response, 0777)
		// Check if the data was updated
		if resp.StatusCode == http.StatusOK && resp.Header.Get("Last-Modified") != "" {
			lastModified, err := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
			if err != nil {
				fmt.Println("Failed to parse Last-Modified header:", err)
				continue
			}

			if lastModified.After(lastUpdateTime) {
				// Data was updated, process the response
				// Here, you can parse the response body or perform any desired actions

				// Update the last update time
				lastUpdateTime = lastModified
			}
		}

		// Close the response body
		resp.Body.Close()

		// Wait for the next poll
		time.Sleep(pollPeriod)
	}
}
