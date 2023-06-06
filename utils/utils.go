package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

func GetPrivateIPAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetWorkingDir get the present working directory
func GetWorkingDir() (string, error) {
	return os.Getwd()
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func GetMacAddresses() ([]string, error) {
	var macAddresses []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// Skip loopback and non-up interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() == nil {
				hwAddr := iface.HardwareAddr
				macAddresses = append(macAddresses, hwAddr.String())
			}
		}
	}

	if len(macAddresses) == 0 {
		return nil, fmt.Errorf("MAC addresses not found")
	}

	return macAddresses, nil
}

func GetGoAgenHash(binaryPath string) (string, error) {
	data, err := ioutil.ReadFile(binaryPath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
