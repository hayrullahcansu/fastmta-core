package conf

//Config is application configs
type Config struct {
	Test        int          `json:"test"`
	IPAddresses []VirtualMta `json:"ip_addressess"`
}

//VirtualMta is ip addresses and define working types
type VirtualMta struct {
	IP       string `json:"ip"`
	Inbound  bool   `json:"inbound"`
	Outbound bool   `json:"outbound"`
}
