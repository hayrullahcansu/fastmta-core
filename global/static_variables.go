package global

import (
	"../conf"
	"github.com/patrickmn/go-cache"
)

var (
	StaticRabbitMqConfig *conf.RabbitMqConfig
	StaticConfig         *conf.Config
	DomainCaches         *cache.Cache
)
