package nmaprunv2

import (
	"encoding/json"
	"encoding/xml"
	"os/exec"

	"github.com/Expand-My-Business/go_windows_agent/utils"
	"github.com/sirupsen/logrus"
)

func RunNmapScan() (*NmapRun, error) {
	// Executing the nmap command from the absloute path
	cmd := exec.Command(`C:\Program Files (x86)\Nmap\nmap.exe`, "-sV", "-oX", "-", "-p-", "localhost")

	// Get the std output
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("Error occured while getting nmap scan output, error: %v", err)
		return nil, err
	}

	var nmapRun NmapRun
	err = xml.Unmarshal(output, &nmapRun)
	if err != nil {
		logrus.Errorf("Error occured while unmarshaling nmap scan output, error: %v", err)
		return nil, err
	}

	return &nmapRun, nil
}

func AddOSTypePortScannedReport(nmap NmapRun) ([]byte, error) {
	for _, p := range nmap.Host.Ports.Port {
		p.Service.Ostype = "Windows"
	}

	privateIP, err := utils.GetPrivateIPAddress()
	if err != nil {
		logrus.Errorf("cannot get private ip: %+v", err)
		return nil, err
	}

	nmapScanReport := NmapScanReport{
		Nmap:   nmap.Host.Ports.Port,
		HostIP: privateIP,
	}

	portScnnedJson, err := json.MarshalIndent(nmapScanReport, "", "\t")
	if err != nil {
		logrus.Errorf("cannot marshal: %+v", err)
		return nil, err
	}

	return portScnnedJson, nil
}

func PortScannedReport() ([]byte, error) {
	nmap, err := RunNmapScan()
	if err != nil {
		logrus.Errorf("cannot get nmap data: %+v", err)
		return nil, err
	}
	return AddOSTypePortScannedReport(*nmap)
}
