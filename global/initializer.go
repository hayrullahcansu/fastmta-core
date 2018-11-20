package global

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/patrickmn/go-cache"

	"../conf"
	OS "../cross"
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
	loadCache()

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

func loadCache() {
	DomainCaches = cache.New(5*time.Minute, 10*time.Minute)
	db, err := entity.GetDbContext()
	if err != nil {
		panic(fmt.Sprintf("Db cant open. %s%s", err, OS.NewLine))
	}
	var domains []entity.Domain
	var count int
	var currentPage int = 0
	var limit int = 50000
	db.Model(&entity.Domain{}).Count(&count)
	for currentPage < count/limit+1 {
		if db.Offset(limit*currentPage).Limit(limit).Model(&entity.Domain).Find(&domains).Error == nil {
			//TODO: burada kaldın foreach loop domains and add to cache
		}
		currentPage++
	}
}
