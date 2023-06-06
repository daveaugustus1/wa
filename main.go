package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Port struct {
	PortID    int    `json:"port_id"`
	State     string `json:"state"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	ExtraInfo string `json:"extra_info"`
}

func main() {
	// Run nmap command
	cmd := exec.Command("nmap", "-p-", "-v", "-A", "127.0.0.1")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute nmap command: %v", err)
	}

	// Parse nmap output
	ports := parseNmapOutput(output)

	// Create JSON file
	file, err := os.Create("nmap_results.json")
	if err != nil {
		log.Fatalf("Failed to create JSON file: %v", err)
	}
	defer file.Close()

	// Write JSON data to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(ports); err != nil {
		log.Fatalf("Failed to encode data to JSON: %v", err)
	}

	fmt.Println("Nmap results have been saved to nmap_results.json")
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
