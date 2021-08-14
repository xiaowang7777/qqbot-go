package goal

import (
	"qqbot-go/config"
	"sync"
)

var (
	ACCESS_TOKEN = ""
	goalConfig   = &config.Config{}

	configSetOne = &sync.Once{}
)

func SetConfig(conf *config.Config) {
	configSetOne.Do(func() {
		goalConfig = conf
	})
}

func GetConfig() *config.Config {
	return goalConfig
}
