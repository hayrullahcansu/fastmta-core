package global

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/patrickmn/go-cache"

	"../conf"
	OS "../cross"
	"../entity"
)

const (
	confFileName   = "app.json"
	dkimFolderPath = "dkim"
)

func Run() {
	confText := readConfigFromFile()
	config := conf.Config{}
	json.Unmarshal([]byte(confText), &config)
	StaticConfig = &config
	StaticRabbitMqConfig = &config.RabbitMq
	dbEnsureCreated()
	go loadDomainCache()

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

func loadDomainCache() {
	DomainCaches = cache.New(5*time.Minute, 10*time.Minute)
	db, err := entity.GetDbContext()
	if err != nil {
		panic(fmt.Sprintf("Db cant open. %s%s", err, OS.NewLine))
	}
	defer db.Close()
	var domains []entity.Domain
	var count int
	var currentPage int = 0
	var limit int = 50000
	db.Model(&entity.Domain{}).Count(&count)
	for currentPage < count/limit+1 {
		if db.Offset(limit*currentPage).Limit(limit).Model(&entity.Domain{}).Find(&domains).Error == nil {
			for _, domain := range domains {
				DomainCaches.Add(domain.DomainName, domain, cache.NoExpiration)
			}
		}
		currentPage++
	}
}

func loadDkimCache() {
	DkimCaches = cache.New(5*time.Minute, 10*time.Minute)
	db, err := entity.GetDbContext()
	if err != nil {
		panic(fmt.Sprintf("Db cant open. %s%s", err, OS.NewLine))
	}
	defer db.Close()
	var dkimmers []entity.Dkimmer
	var count int
	var currentPage int = 0
	var limit int = 1000
	db.Model(&entity.Dkimmer{}).Count(&count)
	for currentPage < count/limit+1 {
		if db.Offset(limit*currentPage).Limit(limit).Model(&entity.Dkimmer{}).Find(&dkimmers).Error == nil {
			for _, dkimmer := range dkimmers {
				DkimCaches.Add(dkimmer.DomainName, dkimmer, cache.NoExpiration)
			}
		}
		currentPage++
	}
	dkimmers = nil
	filepath.Walk(dkimFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		data, err := readFile(path)
		if err != nil {
			return err
		}
		dataString := string(data)
		r := regexp.MustCompile(`(?P<Selector>\w)##(?P<Domain>\w).pem`)
		resultSet := r.FindStringSubmatch(dataString)
		if resultSet != nil {
			dkimmer := &entity.Dkimmer{
				DomainName: resultSet[2],
				Selector:   resultSet[1],
				PrivateKey: dataString,
				Enabled:    true,
			}
			DkimCaches.Add(dkimmer.DomainName, dkimmer, cache.NoExpiration)
		}
		return nil
	})

	//TODO: burada kaldın. GetDkims from folder then push cache
}
func readFile(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)

}
