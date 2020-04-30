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

	"github.com/hayrullahcansu/fastmta-core/util"

	"github.com/hayrullahcansu/fastmta-core/dns"

	"github.com/emersion/go-dkim"
	"github.com/hayrullahcansu/fastmta-core/caching"
	"github.com/hayrullahcansu/fastmta-core/conf"
	"github.com/hayrullahcansu/fastmta-core/constant"
	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/entity"
	"github.com/hayrullahcansu/fastmta-core/global"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/rabbit"
	"github.com/patrickmn/go-cache"
	"github.com/streadway/amqp"
)

const (
	confFileName   = "app.json"
	dkimFolderPath = "dkim"
)

// InitSystem initializes logger and configrations.
// Ensures DB Schema and RabbitMQ queues, exchanges, routing keys...
func InitSystem() {
	logger.Instance()
	config := initConfig(readConfigFromFile())
	global.StaticConfig = config
	global.StaticRabbitMqConfig = &config.RabbitMq
	dbEnsureCreated()
	defineRabbitMqEnvironment()
	go loadDomainCache()
	go loadDkimCache()

}

func initConfig(data string) *conf.Config {
	config := conf.Config{}
	json.Unmarshal([]byte(data), &config)
	if !validationConfig(&config) {
		panic("invalid configration")
	}
	return &config
}

/*

	"mysql", 	"user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local"
	"postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword"
	"sqlite3", 	"/tmp/database_file_name.db"
	"mssql", 	"sqlserver://username:password@localhost:1433?database=dbname"

*/
func validationConfig(config *conf.Config) bool {
	drivers := map[string]byte{
		"sqlite3":  1,
		"mysql":    1,
		"mssql":    1,
		"postgres": 1,
	}
	if _, ok := drivers[config.Database.Driver]; !ok {
		return false
	}
	if util.IsNullOrEmpty(config.Database.Connection) {
		return false
	}
	return true
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
	conn, err := amqp.Dial(rabbit.NewRabbitMqDialString())
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
	err = ch.ExchangeDeclare(constant.InboundExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(constant.InboundStagingExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(constant.OutboundExchange, "direct", true, false, false, false, nil)
	err = ch.ExchangeDeclare(constant.WaitingExchange, "direct", true, false, false, false, nil)

	_, err = ch.QueueDeclare(constant.InboundQueueName, true, false, false, false,
		amqp.Table{
			"x-max-priority": 3,
			"x-queue-mode":   "lazy",
		})
	_, err = ch.QueueDeclare(constant.InboundStagingQueueName, true, false, false, false,
		amqp.Table{
			"x-max-priority": 3,
			"x-queue-mode":   "lazy",
		})
	_, err = ch.QueueDeclare(constant.OutboundNormalQueueName, true, false, false, false,
		amqp.Table{
			"x-max-priority": 3,
			"x-queue-mode":   "lazy",
		})
	_, err = ch.QueueDeclare(constant.OutboundMultipleQueueName, true, false, false, false,
		amqp.Table{
			"x-max-priority": 3,
			"x-queue-mode":   "lazy",
		})
	_, err = ch.QueueDeclare(constant.OutboundWaiting1, true, false, false, false,
		amqp.Table{
			"x-max-priority":            3,
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    constant.WaitingExchange,
			"x-dead-letter-routing-key": constant.RoutingKeyWaiting,
			"x-message-ttl":             1 * 1000,
		})
	_, err = ch.QueueDeclare(constant.OutboundWaiting10, true, false, false, false,
		amqp.Table{
			"x-max-priority":            3,
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    constant.WaitingExchange,
			"x-dead-letter-routing-key": constant.RoutingKeyWaiting,
			"x-message-ttl":             10 * 1000,
		})
	_, err = ch.QueueDeclare(constant.OutboundWaiting60, true, false, false, false,
		amqp.Table{
			"x-max-priority":            3,
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    constant.WaitingExchange,
			"x-dead-letter-routing-key": constant.RoutingKeyWaiting,
			"x-message-ttl":             60 * 1000,
		})
	_, err = ch.QueueDeclare(constant.OutboundWaiting300, true, false, false, false,
		amqp.Table{
			"x-max-priority":            3,
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    constant.WaitingExchange,
			"x-dead-letter-routing-key": constant.RoutingKeyWaiting,
			"x-message-ttl":             300 * 1000,
		})

	// err = ch.QueueBind(constant.InboundQueueName, constant.RoutingKeyInbound, constant.InboundExchange, false, nil)
	fmt.Println(err)
	// err = ch.QueueBind(constant.InboundStagingQueueName, constant.RoutingKeyInboundStaging, constant.InboundStagingExchange, false, nil)
	fmt.Println(err)
	// err = ch.QueueBind(constant.OutboundNormalQueueName, constant.RoutingKeyOutboundNormal, constant.OutboundExchange, false, nil)
	fmt.Println(err)
	// err = ch.QueueBind(constant.OutboundMultipleQueueName, constant.RoutingKeyOutboundMultiple, constant.OutboundExchange, false, nil)
	fmt.Println(err)
	err = ch.QueueBind(constant.OutboundNormalQueueName, constant.RoutingKeyWaiting, constant.WaitingExchange, false, nil)
	fmt.Println(err)
}

func loadDomainCache() {
	domainCacher := caching.InstanceDomain()
	db, err := entity.GetDbContext()
	if err != nil {
		panic(fmt.Sprintf("Db cant open. %s%s", err, cross.NewLine))
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
				if _d, err := dns.NewDomain(domain.DomainName); err == nil {
					domainCacher.AddOrUpdate(domain.DomainName, _d)
				}
			}
		}
		currentPage++
	}
}

func loadDkimCache() {
	dkimCacher := caching.InstanceDkim()
	db, err := entity.GetDbContext()
	if err != nil {
		panic(fmt.Sprintf("Db cant open. %s%s", err, cross.NewLine))
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
