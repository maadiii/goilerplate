package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Name  string
	URL   string
	Test  string
	Admin string
}

var debugConfig = &Config{
	Name:  DEBUG_CONFIG_DATABASE_NAME,
	URL:   DEBUG_CONFIG_DATABASE_URL,
	Test:  DEBUG_CONFIG_DATABASE_TEST,
	Admin: DEBUG_CONFIG_DATABASE_ADMIN,
}

func InitConfig() (*Config, error) {
	var config *Config

	if viper.ConfigFileUsed() == EMPTY {
		config = debugConfig
	} else {
		config = &Config{
			Name:  viper.GetString(CONFIG_DATABASE_NAME),
			URL:   viper.GetString(CONFIG_DATABASE_URL),
			Test:  viper.GetString(CONFIG_DATABASE_TEST),
			Admin: viper.GetString(CONFIG_DATABASE_ADMIN),
		}
	}

	PRE_ERROR := "db/config.go InintConfig(),"
	POST_ERROR := "is not set in config file"
	errMessage := func(msg string) error {
		return fmt.Errorf("%s %s %s", PRE_ERROR, msg, POST_ERROR)
	}

	if config.Name == EMPTY {
		return nil, errMessage(CONFIG_DATABASE_NAME)
	}
	if config.URL == EMPTY {
		return nil, errMessage(CONFIG_DATABASE_URL)
	}
	if config.Test == EMPTY {
		return nil, errMessage(CONFIG_DATABASE_TEST)
	}
	if config.Admin == EMPTY {
		return nil, errMessage(CONFIG_DATABASE_ADMIN)
	}

	return config, nil
}
