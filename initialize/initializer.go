package initialize

import (
	"encoding/json"
	"io/ioutil"

	"../conf"
)

const confFileName = "app.json"

func Run() *conf.Config {
	confText := readConfigFromFile()
	config := conf.Config{}
	json.Unmarshal([]byte(confText), &config)
	return &config
}

func readConfigFromFile() string {
	b, err := ioutil.ReadFile(confFileName)
	if err != nil {
		panic(err)
	}
	return string(b)
}
