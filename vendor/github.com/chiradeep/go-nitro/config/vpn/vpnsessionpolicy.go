package vpn

type Vpnsessionpolicy struct {
	Action  string      `json:"action,omitempty"`
	Builtin interface{} `json:"builtin,omitempty"`
	Name    string      `json:"name,omitempty"`
	Rule    string      `json:"rule,omitempty"`
}
