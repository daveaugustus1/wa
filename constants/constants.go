package constants

const (
	ConfigFilePath = `C:\Program Files\GoAgent\config\config.toml`
)

// Instructions constents
const (
	StartService   = "start"
	StopService    = "stop"
	RestartService = "restart"
	RefreshService = "refresh"
	ScanService    = "scan"
)

const (
	NmapURL                  = "http://13.235.66.99:8011/agent_ports_data"
	NetStatURL               = "http://13.235.66.99:8011/agent_process_data"
	WindowsLogURL            = "http://13.235.66.99:8011/agent_system_logs_data"
	WindowsSystemResourceURL = "http://13.235.66.99:8011/add_agent_logs"
)
