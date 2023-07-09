package instruction

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Expand-My-Business/go_windows_agent/constants"
	"github.com/sirupsen/logrus"
)

func blockDomain(v Instruction) []InstructionResp {
	return BlockDomain(v, []string{v.ServiceName})
}

func BlockDomain(v Instruction, blockedDomains []string) []InstructionResp {
	executionResp := []InstructionResp{}

	hostFile := getHostFilePath()

	// Open the host file in append mode
	file, err := os.OpenFile(hostFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Failed to open host file: %v", err)
		resp := InstructionResp{
			Action:      constants.BlockDomain,
			ID:          v.Id,
			IsExecuted:  false,
			ServiceName: v.ServiceName,
			Status:      "failed",
			Msg:         "",
		}
		executionResp = append(executionResp, resp)
		return executionResp
	}
	defer file.Close()

	// List of blocked domains
	// blockedDomains := []string{
	// 	"www.netflix.com",
	// 	"www.bing.com",
	// }

	// Add blocked domains to the host file
	for _, domain := range blockedDomains {
		_, err = file.WriteString("127.0.0.1 " + domain + "\n")
		if err != nil {
			logrus.Errorf("Failed to write to host file: %v", err)
			resp := InstructionResp{
				Action:      constants.BlockDomain,
				ID:          v.Id,
				IsExecuted:  false,
				ServiceName: v.ServiceName,
				Status:      "failed",
				Msg:         fmt.Sprintf("%v domain couldn't be blocked", v.ServiceName),
			}
			executionResp = append(executionResp, resp)
		} else {
			resp := InstructionResp{
				Action:      constants.BlockDomain,
				ID:          v.Id,
				IsExecuted:  true,
				ServiceName: v.ServiceName,
				Status:      "success",
				Msg:         fmt.Sprintf("%v domain is blocked", v.ServiceName),
			}
			executionResp = append(executionResp, resp)
		}
	}

	logrus.Info("Domains blocked successfully!")

	return executionResp
}

// getHostFilePath returns the path to the host file based on the operating system
func getHostFilePath() string {
	return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
}

func unblockDomain(v Instruction) []InstructionResp {

	executionResp := []InstructionResp{}

	return executionResp
}

// unblockDomains removes the specified domains from the host file
func unblockDomains(v Instruction, domains []string) []InstructionResp {
	hostFile := getHostFilePath()
	executionResp := []InstructionResp{}

	// Read the host file contents
	data, err := os.ReadFile(hostFile)
	if err != nil {
		resp := InstructionResp{
			Action:      constants.BlockDomain,
			ID:          v.Id,
			IsExecuted:  false,
			ServiceName: v.ServiceName,
			Status:      "failed",
			Msg:         fmt.Sprintf("%v domain couldn't be blocked", v.ServiceName),
		}
		executionResp = append(executionResp, resp)
		return executionResp
	}

	// Create a new host file content without the blocked domains
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		shouldRemove := false
		for _, domain := range domains {
			if strings.Contains(line, domain) {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			lines = append(lines, line)
		}
	}

	// Write the updated host file contents
	err = os.WriteFile(hostFile, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		resp := InstructionResp{
			Action:      constants.BlockDomain,
			ID:          v.Id,
			IsExecuted:  false,
			ServiceName: v.ServiceName,
			Status:      "failed",
			Msg:         fmt.Sprintf("%v domain couldn't be blocked", v.ServiceName),
		}
		executionResp = append(executionResp, resp)
		return executionResp
	} else {
		resp := InstructionResp{
			Action:      constants.BlockDomain,
			ID:          v.Id,
			IsExecuted:  true,
			ServiceName: v.ServiceName,
			Status:      "success",
			Msg:         fmt.Sprintf("%v domain is blocked", v.ServiceName),
		}
		executionResp = append(executionResp, resp)
		return executionResp
	}

}
