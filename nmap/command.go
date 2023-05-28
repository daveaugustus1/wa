package nmap

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Port struct {
	PortID    int    `json:"port_id"`
	State     string `json:"state"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	ExtraInfo string `json:"extra_info"`
}

func NmapDataCmd() ([]byte, error) {
	// Run nmap command
	cmd := exec.Command("nmap", "-p-", "-v", "-A", "127.0.0.1")
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("Failed to execute nmap command: %v", err)
		return nil, err
	}

	// Parse nmap output
	ports := parseNmapOutput(output)

	data, err := json.MarshalIndent(ports, "", "\t")
	if err != nil {
		logrus.Errorf("Failed to marshal nmap data: %v", err)
		return nil, err
	}

	logrus.Infof("Nmap results have been saved to nmap_results.json")
	return data, err
}

func parseNmapOutput(output []byte) []Port {
	var ports []Port

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Check if the line contains port information
		if strings.HasPrefix(line, "PORT") || strings.HasPrefix(line, "Nmap scan report") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 3 {
			portID, _ := strconv.Atoi(strings.Split(fields[0], "/")[0])
			state := fields[1]
			service := fields[2]
			version := ""
			extraInfo := ""
			if len(fields) >= 4 {
				version = fields[3]
			}
			if len(fields) >= 5 {
				extraInfo = strings.Join(fields[4:], " ")
			}

			port := Port{
				PortID:    portID,
				State:     state,
				Service:   service,
				Version:   version,
				ExtraInfo: extraInfo,
			}

			ports = append(ports, port)
		}
	}

	return ports
}
