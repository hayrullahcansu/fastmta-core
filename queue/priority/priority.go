package priority

//Priority for messages in RabbitMQ queues.
type Priority uint8

const (
	//LOW is equals 0 priority
	LOW Priority = iota + 0
	//NORMAL is equals 1 priority
	NORMAL
	//HIGH is equals 2 prioprity
	HIGH
)
