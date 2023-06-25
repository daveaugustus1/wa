package instruction

type InstructionSet struct {
	CompanyCode string        `json:"company_code"`
	AgentIP     string        `json:"agent_ip"`
	MacAddress  string        `json:"mac_address"`
	Instruction []Instruction `json:"instruction"`
}

type Instruction struct {
	ServiceName string `json:"service_name"`
	Action      string `json:"action"`
	IsExecuted  bool   `json:"is_executed"`
}
