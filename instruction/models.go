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

type InstructionSetResp struct {
	AgentIP          string            `json:"agent_ip"`
	CompanyCode      string            `json:"company_code"`
	InstructionResps []InstructionResp `json:"instruction"`
	MacAddress       string            `json:"mac_address"`
}

type InstructionResp struct {
	Action      string `json:"action"`
	ID          string `json:"id"`
	IsExecuted  bool   `json:"is_executed"`
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
	Msg         string `json:"msg"`
}
