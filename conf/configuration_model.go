package conf

//Config is application configs
type Config struct {
	Test        int              `json:"test"`
	IPAddresses []VirtualMta     `json:"ip_addressess"`
	RabbitMq    []RabbitMqConfig `json:"rabbitmq"`
}

//VirtualMta is ip addresses and define working types
type VirtualMta struct {
	IP       string `json:"ip"`
	Inbound  bool   `json:"inbound"`
	Outbound bool   `json:"outbound"`
}

//RabbitMqConfig is rabbitmq configuration to connect
type RabbitMqConfig struct {
	HostName     string `json:"hostname"`
	VirtualHost  string `json:"virtual_host"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Port         int    `json:"port"`
	ExchangeName string `json:"exchange_name"`
}
