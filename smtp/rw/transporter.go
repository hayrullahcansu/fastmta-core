package rw

type Transporter interface {
	ReadData() (string, error)
	ReadAllLine() (string, error)
	ReadAll() (string, error)
	WriteLine(data string) error
}
