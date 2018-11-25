package global

import (
	"../conf"
	"github.com/patrickmn/go-cache"
)

//TODO: instances define onb itself
var (
	StaticRabbitMqConfig *conf.RabbitMqConfig
	StaticConfig         *conf.Config
	DomainCaches         *cache.Cache
	DkimCaches           *cache.Cache
)
