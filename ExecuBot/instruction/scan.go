package instruction

import (
	"github.com/Expand-My-Business/go_windows_agent/constants"
	"github.com/Expand-My-Business/go_windows_agent/goagent/config"
	"github.com/Expand-My-Business/go_windows_agent/goagent/netstat"
	"github.com/Expand-My-Business/go_windows_agent/goagent/nmaprunv2"
	"github.com/Expand-My-Business/go_windows_agent/goagent/windowslogs"
	"github.com/Expand-My-Business/go_windows_agent/utils"
)

func scanOperation(v Instruction) []InstructionResp {
	cfg := config.GetConfigInstance()

	executionResp := []InstructionResp{}

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
			executionResp = append(executionResp, executionRes)
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
				executionResp = append(executionResp, executionRes)
			} else {
				executionRes := InstructionResp{
					ID:          v.Id,
					Action:      constants.ScanService,
					IsExecuted:  true,
					Msg:         "The service scanned successfully",
					ServiceName: v.ServiceName,
					Status:      "passed",
				}
				executionResp = append(executionResp, executionRes)
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
			executionResp = append(executionResp, executionRes)
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
				executionResp = append(executionResp, executionRes)
			} else {
				executionRes := InstructionResp{
					ID:          v.Id,
					Action:      constants.ScanService,
					IsExecuted:  true,
					Msg:         "The service scanned successfully",
					ServiceName: v.ServiceName,
					Status:      "passed",
				}
				executionResp = append(executionResp, executionRes)
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
			executionResp = append(executionResp, executionRes)
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
				executionResp = append(executionResp, executionRes)
			} else {
				executionRes := InstructionResp{
					ID:          v.Id,
					Action:      constants.ScanService,
					IsExecuted:  true,
					Msg:         "The service scanned successfully",
					ServiceName: v.ServiceName,
					Status:      "passed",
				}
				executionResp = append(executionResp, executionRes)
			}
		}
	}
	return executionResp
}
