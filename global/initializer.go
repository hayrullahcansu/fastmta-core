package global

import (
	"encoding/json"
	"io/ioutil"

	"../conf"
	"../entity"
)

const confFileName = "app.json"

func Run() {
	confText := readConfigFromFile()
	config := conf.Config{}
	json.Unmarshal([]byte(confText), &config)
	StaticConfig = &config
	StaticRabbitMqConfig = &config.RabbitMq
	dbEnsureCreated()
}

func readConfigFromFile() string {
	b, err := ioutil.ReadFile(confFileName)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func dbEnsureCreated() {
	db, _ := entity.GetDbContext()
	db.Create(&entity.Domain{
		DomainName: "gmail.com",
		MXRecords: []entity.MXRecord{
			entity.MXRecord{
				Host:           "mx1.gmail.com",
				Pref:           123,
				BaseDomainName: "gmail.com",
			},
			entity.MXRecord{
				Host:           "mx2.gmail.com",
				Pref:           12,
				BaseDomainName: "gmail.com",
			},
		},
	})
	db.Close()
}
