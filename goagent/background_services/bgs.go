package background_services

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type Service struct {
	Name        string
	Status      string
	DisplayName string
}

func ServicesCanBeStopped() []Service {
	// Command to execute
	cmd := exec.Command("powershell", "-Command", "Get-Service | Where-Object { $_.CanStop }")

	// Execute the command and capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute command: %v\n%s", err, stderr.String())
	}

	// Process the output
	output := stdout.String()
	log.Println("Command output:\n", output)

	// Parse the output and create service structs
	services := parseServiceOutput(output)

	return services
}

// Parses the output of Get-Service command and returns a slice of Service structs
func parseServiceOutput(output string) []Service {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	services := make([]Service, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			service := Service{
				Status:      fields[0],
				Name:        fields[1],
				DisplayName: strings.Join(fields[3:], " "),
			}
			services = append(services, service)
		}
	}

	return services
}
