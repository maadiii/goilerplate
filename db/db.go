package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Session struct {
	*gorm.DB
}

func NewSession(config *Config) (*Session, error) {
	db, err := gorm.Open(config.Name, config.URL)
	if err != nil {
		return nil,
			fmt.Errorf("db/db.go New(), Unable connect to database. %v", err)
	}

	return &Session{db}, nil
}
