package datastore

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var debugConfig = &Config{
	Driver: DEBUG_CONFIG_DATABASE_DRIVER,
	URL:    DEBUG_CONFIG_DATABASE_URL,
	Test:   DEBUG_CONFIG_DATABASE_TEST,
	Admin:  DEBUG_CONFIG_DATABASE_ADMIN,
}

type Config struct {
	Driver string
	URL    string
	Test   string
	Admin  string
}

type Session struct {
	*gorm.DB
}

func InitConfig() (*Config, error) {
	var config *Config

	if viper.ConfigFileUsed() == EMPTY {
		config = debugConfig
	} else {
		config = &Config{
			Driver: viper.GetString(CONFIG_DATABASE_DRIVER),
			URL:    viper.GetString(CONFIG_DATABASE_URL),
			Test:   viper.GetString(CONFIG_DATABASE_TEST),
			Admin:  viper.GetString(CONFIG_DATABASE_ADMIN),
		}
	}

	PRE_ERROR := "db/config.go InintConfig(),"
	POST_ERROR := "is not set in config file"
	errMessage := func(msg string) error {
		return fmt.Errorf("%s %s %s", PRE_ERROR, msg, POST_ERROR)
	}

	if config.Driver == EMPTY {
		return nil, errMessage(CONFIG_DATABASE_DRIVER)
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

func NewSession(config *Config) (*Session, error) {
	db, err := gorm.Open(config.Driver, config.URL)
	if err != nil {
		return nil,
			fmt.Errorf("db/db.go New(), Unable connect to database. %v", err)
	}

	return &Session{db}, nil
}

func NewTestSession(c *Config) (*Session, error) {
	db, err := gorm.Open(c.Driver, c.Test)
	if err != nil {
		return nil, err
	}

	return &Session{db}, nil
}

const (
	EMPTY = ""

	// debug config values
	DEBUG_CONFIG_DATABASE_DRIVER = "postgres"
	DEBUG_CONFIG_DATABASE_URL    = "host=localhost port=5432 " +
		"user=goilerplate password=goilerplate dbname=goilerplate"
	DEBUG_CONFIG_DATABASE_TEST = "host=localhost port=5432 " +
		"user=goilerplate password=goilerplate dbname=goilerplate_test"
	DEBUG_CONFIG_DATABASE_ADMIN = "host=localhost port=5432 user=postgres " +
		"password=postgres dbname=postgres"

	// config names
	CONFIG_DATABASE_DRIVER = "database.driver"
	CONFIG_DATABASE_URL    = "database.url"
	CONFIG_DATABASE_TEST   = "database.test"
	CONFIG_DATABASE_ADMIN  = "database.admin"
)
