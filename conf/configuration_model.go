package conf

//Config is application configs
type Config struct {
	Test        int            `json:"test"`
	IPAddresses []VirtualMta   `json:"ip_addressess"`
	RabbitMq    RabbitMqConfig `json:"rabbitmq"`
}

//VirtualMta is ip addresses and define working types
type VirtualMta struct {
	IP       string `json:"ip"`
	Inbound  bool   `json:"inbound"`
	Outbound bool   `json:"outbound"`
}

//RabbitMqConfig is rabbitmq configuration to connect
type RabbitMqConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	VirtualHost  string `json:"virtual_host"`
	ExchangeName string `json:"exchange_name"`
}
