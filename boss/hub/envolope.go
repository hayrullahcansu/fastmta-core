package hub

type Envelope struct {
	Client  string
	Command string
	Target  string
	Message interface{} `json:"omitempty"`
}
