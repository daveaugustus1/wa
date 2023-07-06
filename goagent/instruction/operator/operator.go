package operator

import (
	"fmt"
	"os/exec"
	"strings"
)

// func main() {
// 	serviceName := "GoAgent"

// 	// // Start the service
// 	// err := startService(serviceName)
// 	// if err != nil {
// 	// 	fmt.Println("Failed to start service:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("Service started successfully.")

// 	// // Stop the service
// 	// err = stopService(serviceName)
// 	// if err != nil {
// 	// 	fmt.Println("Failed to stop service:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("Service stopped successfully.")

// 	// Restart the service
// 	err := RestartService(serviceName)
// 	if err != nil {
// 		fmt.Println("Failed to restart service:", err)
// 		return
// 	}
// 	fmt.Println("Service restarted successfully.")
// }

func StartService(serviceName string) error {
	command := fmt.Sprintf("Start-Service -Name %s", serviceName)
	output, err := runPowerShellCommand(command)
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	// Check the output for any error messages or confirmations if needed
	// For example, you can parse the output to check if the service was successfully started
	_ = output

	return nil
}

func StopService(serviceName string) error {
	command := fmt.Sprintf("Stop-Service -Name %s", serviceName)
	output, err := runPowerShellCommand(command)
	if err != nil {
		return fmt.Errorf("failed to stop service: %v", err)
	}

	// Check the output for any error messages or confirmations if needed
	// For example, you can parse the output to check if the service was successfully stopped
	_ = output

	return nil
}

func RestartService(serviceName string) error {
	command := fmt.Sprintf("Restart-Service -Name %s", serviceName)
	output, err := runPowerShellCommand(command)
	if err != nil {
		return fmt.Errorf("failed to restart service: %v", err)
	}

	// Check the output for any error messages or confirmations if needed
	// For example, you can parse the output to check if the service was successfully restarted
	_ = output

	return nil
}

func runPowerShellCommand(command string) (string, error) {
	psCmd := exec.Command("PowerShell", "-Command", command)
	output, err := psCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run PowerShell command: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}
