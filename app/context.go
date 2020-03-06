package app

import (
	"goilerplate/db"
	"goilerplate/domain/models"

	"github.com/sirupsen/logrus"
)

type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	DBSession     *db.Session
	User          *models.User
	Model         []byte
}

func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := ctx
	ret.Logger = logger
	return ret
}

func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := ctx
	ret.RemoteAddress = address
	return ret
}

func (ctx *Context) WithUser(user *models.User) *Context {
	ret := ctx
	ret.User = user
	return ret
}

func (ctx *Context) WithModel(model []byte) *Context {
	ret := ctx
	ret.Model = model
	return ret
}
