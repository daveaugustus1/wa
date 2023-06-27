package instruction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Expand-My-Business/go_windows_agent/config"
	"github.com/Expand-My-Business/go_windows_agent/constants"
	"github.com/Expand-My-Business/go_windows_agent/instruction/operator"
	"github.com/Expand-My-Business/go_windows_agent/utils"
	"github.com/sirupsen/logrus"
)

const (
	pollPeriod = 10 * time.Second // Polling interval in seconds
)

func GetInstructions() {
	privateIP, err := utils.GetPrivateIPAddress()
	if err != nil {
		logrus.Errorf("Cannot get the private ip, error: %+v", err)
		return
	}

	// lastUpdateTime := time.Time{}
	cfg := config.GetConfigInstance()
	for {
		// Prepare the request
		url := fmt.Sprintf(cfg.InstructionEndpoint, cfg.Port, privateIP)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			logrus.Errorf("Failed to create request, error: %v", err)
			continue
		}

		// Set headers
		req.Header.Set("company-code", cfg.CompanyCode)

		// Send the request
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Request failed, error: %v", err)
			continue
		}

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Cannot read response body, error: %+v", err)
			return
		}

		ins := InstructionSet{}
		json.Unmarshal(response, &ins)
		fmt.Printf("ins: %+v\n", ins)

		if ins.CompanyCode != cfg.CompanyCode {
			logrus.Errorf("deployed company code %s didn't match with resp code %s", cfg.CompanyCode, ins.CompanyCode)
			return
		}

		// Check if the data was updated
		if resp.StatusCode == http.StatusOK {
			for _, v := range ins.Instruction {
				switch v.Action {
				case constants.StartService:
					if err := operator.StartService(v.ServiceName); err != nil {
						logrus.Errorf("cannot start the service, error: %v", err)
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service started successfully")
					}
				case constants.StopService:
					if err := operator.StopService(v.ServiceName); err != nil {
						logrus.Errorf("cannot stop the service, error: %v", err)
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service stoped successfully")
					}
				case constants.RestartService:
					if err := operator.RestartService(v.ServiceName); err != nil {
						logrus.Errorf("cannot restart the service, error: %v", err)
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service restarted successfully")
					}
				case constants.RefreshService:
					if err := operator.RestartService(v.ServiceName); err != nil {
						logrus.Errorf("cannot restart the service, error: %v", err)
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service refreshed successfully")
					}
				default:
					RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "Well, the action isn't supported by the agent!")
				}
			}
		}

		// Close the response body
		resp.Body.Close()

		// Wait for the next poll
		time.Sleep(pollPeriod)
	}
}

type Payload struct {
	Servicename string `json:"service_name"`
	Message     string `json:"message"`
}

func RespondExecutionDetails(endpoint, serviceName, companyCode, msg string) {
	// Prepare the payload data
	payload := Payload{
		Servicename: serviceName,
		Message:     msg,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal payload: %v", err)
		return
	}

	// Send POST request
	logrus.Infof("Instruction resp url: %v", endpoint)
	// url := "https://api.example.com/endpoint" // Replace with your API endpoint URL
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logrus.Errorf("Failed to create POST request: %v", err)
		return
	}

	// Set request headers (if needed)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("company-code", companyCode)

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("POST request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("API returned non-OK status: %v", resp.Status)
		return
	}

	logrus.Infof("POST request completed successfully.")
}
