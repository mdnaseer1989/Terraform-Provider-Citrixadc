package vpn

type Vpnicaconnection struct {
	Destip   string `json:"destip,omitempty"`
	Destport int    `json:"destport,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Peid     int    `json:"peid,omitempty"`
	Srcip    string `json:"srcip,omitempty"`
	Srcport  int    `json:"srcport,omitempty"`
	Username string `json:"username,omitempty"`
}
