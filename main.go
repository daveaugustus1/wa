package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	BlockDomain()
}

func BlockDomain() {
	hostFile := getHostFilePath()

	// Open the host file in append mode
	file, err := os.OpenFile(hostFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open host file:", err)
		return
	}
	defer file.Close()

	// List of blocked domains
	blockedDomains := []string{
		"www.netflix.com",
		"www.bing.com",
	}

	// Add blocked domains to the host file
	for _, domain := range blockedDomains {
		_, err = file.WriteString("127.0.0.1 " + domain + "\n")
		if err != nil {
			fmt.Println("Failed to write to host file:", err)
			return
		}
	}

	fmt.Println("Domains blocked successfully!")
}

// getHostFilePath returns the path to the host file based on the operating system
func getHostFilePath() string {
	return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
}

func unblockDOmain() {
	hostFile := getHostFilePath()

	// Remove blocked domains from the host file
	unblockDomains(hostFile, []string{
		"www.netflix.com",
		"www.bing.com",
	})

	fmt.Println("Domains unblocked successfully!")
}

// unblockDomains removes the specified domains from the host file
func unblockDomains(hostFile string, domains []string) {
	// Read the host file contents
	data, err := os.ReadFile(hostFile)
	if err != nil {
		fmt.Println("Failed to read host file:", err)
		return
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
		fmt.Println("Failed to write to host file:", err)
		return
	}
}

// // getHostFilePath returns the path to the host file based on the operating system
// func getHostFilePath() string {
// 	return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
// }
