package conf

//Config is application configs
type Config struct {
	IPAddresses []VirtualMta   `json:"ip_addressess"`
	Ports       []int          `json:"ports"`
	RabbitMq    RabbitMqConfig `json:"rabbitmq"`
	Database    DatabaseConfig `json:"database"`
}

//VirtualMta is ip addresses and define working types
type VirtualMta struct {
	IP       string `json:"ip"`
	HostName string `json:"hostname"`
	Inbound  bool   `json:"inbound"`
	Outbound bool   `json:"outbound"`
	GroupID  int    `json:"group_id"`
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

//RabbitMqConfig is rabbitmq configuration to connect
type DatabaseConfig struct {
	Driver     string `json:"driver"`
	Connection string `json:"connection"`
}
