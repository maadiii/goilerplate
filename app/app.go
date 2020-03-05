package app

import (
	"goldfish/db"

	"github.com/sirupsen/logrus"
)

var (
	Version = VERSION
	Name    = NAME
)

type App struct {
	Config    *Config
	DBSession *db.Session
}

func (a *App) NewContext() *Context {
	return &Context{
		Logger:    logrus.StandardLogger(),
		DBSession: a.DBSession,
	}
}

func New() (app *App, err error) {
	app = &App{}
	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	app.DBSession, err = db.NewSession(dbConfig)
	if err != nil {
		return nil, err
	}

	return app, err
}

func (a *App) Close() error {
	return a.DBSession.Close()
}
