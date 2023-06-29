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
	"github.com/Expand-My-Business/go_windows_agent/netstat"
	"github.com/Expand-My-Business/go_windows_agent/nmaprunv2"
	"github.com/Expand-My-Business/go_windows_agent/utils"
	"github.com/Expand-My-Business/go_windows_agent/windowslogs"
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

		ip, _ := utils.GetPrivateIPAddress()
		macAddr, _ := utils.GetMacAddresses()
		// Create response struct
		executionResp := InstructionSetResp{
			AgentIP:     ip,
			MacAddress:  macAddr[0],
			CompanyCode: cfg.CompanyCode,
		}

		// Check if the data was updated
		if resp.StatusCode == http.StatusOK {
			for _, v := range ins.Instruction {
				switch v.Action {
				case constants.StartService:
					if err := operator.StartService(v.ServiceName); err != nil {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.StartService,
							IsExecuted:  true,
							Msg:         err.Error(),
							ServiceName: v.ServiceName,
							Status:      "failed",
						}
						logrus.Errorf("cannot start the service %s, error: %v", v.ServiceName, err)
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.StartService,
							IsExecuted:  true,
							Msg:         "The service started successfully",
							ServiceName: v.ServiceName,
							Status:      "passed",
						}
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service started successfully")
					}
				case constants.StopService:
					if err := operator.StopService(v.ServiceName); err != nil {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.StopService,
							IsExecuted:  true,
							Msg:         err.Error(),
							ServiceName: v.ServiceName,
							Status:      "failed",
						}
						logrus.Errorf("cannot stop the service %s, error: %v", v.ServiceName, err)
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.StopService,
							IsExecuted:  true,
							Msg:         "The service stoped successfully",
							ServiceName: v.ServiceName,
							Status:      "passed",
						}
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service stoped successfully")
					}
				case constants.RestartService:
					if err := operator.RestartService(v.ServiceName); err != nil {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.RestartService,
							IsExecuted:  true,
							Msg:         err.Error(),
							ServiceName: v.ServiceName,
							Status:      "failed",
						}
						logrus.Errorf("cannot restart the service %s, error: %v", v.ServiceName, err)
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, err.Error())
					} else {
						executionRes := InstructionResp{
							ID:          v.Id,
							Action:      constants.RestartService,
							IsExecuted:  true,
							Msg:         "The service restarted successfully",
							ServiceName: v.ServiceName,
							Status:      "passed",
						}
						executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						// RespondExecutionDetails(cfg.InstructionRespEndpoint, v.ServiceName, ins.CompanyCode, "The service restarted successfully")
					}
				case constants.ScanService:
					if v.ServiceName == "nmap" {
						nmapXbyte, err := nmaprunv2.PortScannedReport()
						if err != nil {
							executionRes := InstructionResp{
								ID:          v.Id,
								Action:      constants.ScanService,
								IsExecuted:  true,
								Msg:         err.Error(),
								ServiceName: v.ServiceName,
								Status:      "failed",
							}
							executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						} else {
							// TODO: trigger response concurrently
							if err := utils.SendStringToAPI(constants.NmapURL, string(nmapXbyte), cfg.CompanyCode); err != nil {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         err.Error(),
									ServiceName: v.ServiceName,
									Status:      "failed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							} else {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         "The service scanned successfully",
									ServiceName: v.ServiceName,
									Status:      "passed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							}
						}
					} else if v.ServiceName == "netstat" {
						netXbyte, err := netstat.GetNetStats()
						if err != nil {
							executionRes := InstructionResp{
								ID:          v.Id,
								Action:      constants.ScanService,
								IsExecuted:  true,
								Msg:         err.Error(),
								ServiceName: v.ServiceName,
								Status:      "failed",
							}
							executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						} else {
							// TODO: trigger response concurrently
							if err := utils.SendStringToAPI(constants.NetStatURL, string(netXbyte), cfg.CompanyCode); err != nil {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         err.Error(),
									ServiceName: v.ServiceName,
									Status:      "failed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							} else {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         "The service scanned successfully",
									ServiceName: v.ServiceName,
									Status:      "passed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							}
						}
					} else if v.ServiceName == "system-log" {
						windoesLogs, err := windowslogs.GetSystemLogs()
						if err != nil {
							executionRes := InstructionResp{
								ID:          v.Id,
								Action:      constants.ScanService,
								IsExecuted:  true,
								Msg:         err.Error(),
								ServiceName: v.ServiceName,
								Status:      "failed",
							}
							executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
						} else {
							// TODO: trigger response concurrently
							if err := utils.SendStringToAPI(constants.WindowsLogURL, string(windoesLogs), cfg.CompanyCode); err != nil {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         err.Error(),
									ServiceName: v.ServiceName,
									Status:      "failed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							} else {
								executionRes := InstructionResp{
									ID:          v.Id,
									Action:      constants.ScanService,
									IsExecuted:  true,
									Msg:         "The service scanned successfully",
									ServiceName: v.ServiceName,
									Status:      "passed",
								}
								executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
							}
						}
					}
				default:
					executionRes := InstructionResp{
						ID:          v.Id,
						Action:      v.Action,
						IsExecuted:  true,
						Msg:         "Well, the action isn't supported by the agent!",
						ServiceName: v.ServiceName,
						Status:      "failed",
					}
					executionResp.InstructionResps = append(executionResp.InstructionResps, executionRes)
				}
			}
		}

		RespondExecutionDetails(cfg.InstructionRespEndpoint, cfg.CompanyCode, executionResp)
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

func RespondExecutionDetails(endpoint, companyCode string, payload InstructionSetResp) {
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
