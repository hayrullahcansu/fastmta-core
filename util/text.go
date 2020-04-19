package util

import "strings"

func readToClRf() string {
	//var b strings.Builder
	//reader := bufio.NewReader(io.Reader)
	return ""
}

//IsNullOrEmpty checks input value is empty which is string
func IsNullOrEmpty(data string) bool {
	return len(strings.TrimSpace(data)) <= 0
}

//IsNullOrEmptyB checks input value is empty which is bytes
func IsNullOrEmptyB(data []byte) bool {
	return len(strings.TrimSpace(string(data))) <= 0
}
