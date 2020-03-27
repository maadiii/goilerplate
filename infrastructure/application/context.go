package application

import (
	"encoding/json"
	"fmt"
	"goilerplate/domain/models"
	"goilerplate/infrastructure/datastore"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Context struct {
	Request       *http.Request
	Response      http.ResponseWriter
	Logger        logrus.FieldLogger
	RemoteAddress string
	DBSession     *datastore.Session
	User          *models.User
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

func (ctx *Context) WithRequest(request *http.Request) *Context {
	ret := ctx
	ret.Request = request
	return ret
}

func (ctx *Context) WithResponseWriter(responseWriter http.ResponseWriter) *Context {
	ret := ctx
	ret.Response = responseWriter
	return ret
}

func (ctx *Context) DecodeModel(v interface{}) {
	if err := json.NewDecoder(ctx.Request.Body).Decode(v); err != nil {
		panic(fmt.Sprintf("failed to decode model because %v", err))
	}
}

func (ctx *Context) Json(v interface{}) {
	if err := json.NewEncoder(ctx.Response).Encode(v); err != nil {
		panic(fmt.Sprintf("failed to encode json to response %v", err))
	}
}
