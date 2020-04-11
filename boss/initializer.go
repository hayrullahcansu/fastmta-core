package boss

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/emersion/go-dkim"
	"github.com/hayrullahcansu/fastmta-core/caching"
	"github.com/hayrullahcansu/fastmta-core/conf"
	OS "github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/global"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/queue"
	"github.com/patrickmn/go-cache"
	"github.com/streadway/amqp"
)

const (
	confFileName   = "app.json"
	dkimFolderPath = "dkim"
)

func InitSystem() {
	logger.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	confText := readConfigFromFile()
	config := conf.Config{}
	json.Unmarshal([]byte(confText), &config)
	global.StaticConfig = &config
	global.StaticRabbitMqConfig = &config.RabbitMq
	dbEnsureCreated()
	defineRabbitMqEnvironment()
	go loadDomainCache()
	go loadDkimCache()

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

func defineRabbitMqEnvironment() {
	conn, err := amqp.Dial(queue.NewRabbitMqDialString())
	if err != nil {
		//FIXME: Send signal to main process to kill
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		//FIXME: Send signal to main process to kill
		panic(err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(queue.InboundExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(queue.InboundStagingExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(queue.OutboundExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(queue.OutboundExchange, "direct", true, false, false, false, nil)

	_, err = ch.QueueDeclare(queue.InboundQueueName, true, false, false, false, nil)
	_, err = ch.QueueDeclare(queue.InboundStagingQueueName, true, false, false, false, nil)
	_, err = ch.QueueDeclare(queue.OutboundNormalQueueName, true, false, false, false, nil)
	_, err = ch.QueueDeclare(queue.OutboundMultipleQueueName, true, false, false, false, nil)

	err = ch.QueueBind(queue.InboundQueueName, queue.RoutingKeyInbound, queue.InboundExchange, false, nil)
	fmt.Println(err)
	err = ch.QueueBind(queue.InboundStagingQueueName, queue.RoutingKeyInboundStaging, queue.InboundStagingExchange, false, nil)
	fmt.Println(err)
	err = ch.QueueBind(queue.OutboundNormalQueueName, queue.RoutingKeyOutboundNormal, queue.OutboundExchange, false, nil)
	fmt.Println(err)
	err = ch.QueueBind(queue.OutboundMultipleQueueName, queue.RoutingKeyOutboundMultiple, queue.OutboundExchange, false, nil)
	fmt.Println(err)

}

func loadDomainCache() {
	domainCacher := caching.InstanceDomain()
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
				domainCacher.C.Add(domain.DomainName, domain, cache.NoExpiration)
			}
		}
		currentPage++
	}
}

func loadDkimCache() {
	dkimCacher := caching.InstanceDkim()
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

				options, err := getDkimOption(dkimmer.DomainName, dkimmer.Selector, dkimmer.PrivateKey)
				if err == nil {
					dkimmer.Options = options
					dkimCacher.C.Add(dkimmer.DomainName, dkimmer, cache.NoExpiration)
				}
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
			options, err := getDkimOption(dkimmer.DomainName, dkimmer.Selector, dkimmer.PrivateKey)
			if err == nil {
				dkimmer.Options = options
				dkimCacher.C.Add(dkimmer.DomainName, dkimmer, cache.NoExpiration)
			}
			dkimCacher.C.Add(dkimmer.DomainName, dkimmer, cache.NoExpiration)
		}
		return nil
	})
}
func readFile(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		//log.Fatalf("open file error: %v", err)
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)

}

func getDkimOption(domain string, selector string, privateKey string) (*dkim.SignOptions, error) {
	block, _ := pem.Decode([]byte(privateKey))
	signer, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		options := &dkim.SignOptions{
			Domain:   domain,
			Selector: selector,
			Signer:   signer,
		}
		return options, nil
	}
	return nil, err
}
