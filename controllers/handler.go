package controllers

import (
	"fmt"
	"goilerplate/app"
	"goilerplate/domain/models"
	"goilerplate/domain/services"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

type Config struct {
	HTTPSPort  int
	HTTPPort   int
	DomainName string
	ProxyCount int
}

var debugConfig = &Config{
	HTTPSPort:  8443,
	HTTPPort:   8000,
	DomainName: "127.0.0.1",
}

func InitConfig() (*Config, error) {
	var config *Config
	PRE_ERROR := "controllers/handler.go InitConfig()"

	if viper.ConfigFileUsed() == EMPTY {
		config = debugConfig
	} else {
		config = &Config{
			HTTPSPort:  viper.GetInt(HTTPS_PORT),
			HTTPPort:   viper.GetInt(HTTP_PORT),
			DomainName: viper.GetString(DOMAIN_NAME),
			ProxyCount: viper.GetInt(PROXY_COUNT),
		}
	}
	if config.HTTPPort == 0 {
		return nil, fmt.Errorf("%s %s", PRE_ERROR, HTTP_ERROR)
	}
	if config.HTTPSPort == 0 {
		return nil, fmt.Errorf("%s %s", PRE_ERROR, HTTPS_ERROR)
	}

	return config, nil
}

type Controller struct {
	App      *app.App
	Config   *Config
	Request  *http.Request
	Response http.ResponseWriter
}

func New(a *app.App) (controller *Controller, err error) {
	controller = &Controller{App: a}
	controller.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func (ctrl *Controller) IPAddressForRequest(r *http.Request) string {
	addr := r.RemoteAddr
	if ctrl.Config.ProxyCount > 0 {
		h := r.Header.Get(X_FORWARDED_FOR)
		if h != EMPTY {
			clients := strings.Split(h, COMMA)
			if ctrl.Config.ProxyCount > len(clients) {
				addr = clients[0]
			} else {
				addr = clients[len(clients)-ctrl.Config.ProxyCount]
			}
		}
	}
	return strings.Split(strings.TrimSpace(addr), COLON)[0]
}

type action func(*app.Context) error

func (ctrl *Controller) Authorize(f action, role string) action {
	return func(ctx *app.Context) error {
		tokenString, err := ctrl.App.ReadCookie(ctrl.Request, ACCESS_COOKIE)
		if err != nil {
			return app.NewErrUnauthorized()
		}

		claims := &app.Claims{}
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return ctrl.App.Config.JWT.Secret, nil
		}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			keyFunc,
		)
		if err != nil {
			return app.NewErrUnauthorized()
		}

		if !token.Valid {
			refreshTokenString, err := ctrl.App.
				ReadCookie(ctrl.Request, REFRESH_COOKIE)
			if err != nil {
				return app.NewErrUnauthorized()
			}

			refreshClaims := &app.Claims{}
			refreshToken, err := jwt.ParseWithClaims(
				refreshTokenString,
				refreshClaims,
				keyFunc,
			)
			if err != nil {
				return app.NewErrUnauthorized()
			}

			if !refreshToken.Valid {
				return app.NewErrUnauthorized()
			}

			if refreshClaims.ID != claims.ID {
				return app.NewErrUnauthorized()
			}

			user := models.User{ID: refreshClaims.ID}
			err = services.NewUserService(ctx.DBSession).
				GetUserWithGroupAndRole(&user)
			if err != nil {
				return app.NewErrUnauthorized()
			}

			roles := make([]string, len(user.Group.Roles))
			for i, role := range user.Group.Roles {
				roles[i] = role.EnName
			}

			jwt, err := ctrl.App.CreateJWT(
				user.ID,
				user.FirstName,
				user.LastName,
				false,
				roles...,
			)
			if err != nil {
				return err
			}

			if err := ctrl.App.SetCookie(
				ACCESS_COOKIE,
				jwt,
				false,
				ctrl.Response,
			); err != nil {
				return err
			}

			claims.Roles = roles
		}

		user := &models.User{
			ID:        claims.ID,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
		}
		ctx = ctx.WithUser(user)

		doNext := false
		if role != EMPTY {
			for _, r := range claims.Roles {
				if r == role {
					doNext = true
					break
				}
			}
			if !doNext {
				return app.NewErrForbidden()
			}
		}

		return f(ctx)
	}
}

func (ctrl *Controller) HandleView(f action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		r.ParseForm()
		defer r.Body.Close()

		ctx := ctrl.App.NewContext().
			WithRemoteAddress(ctrl.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger)

		tempModel := `{`
		for key, value := range r.Form {
			tempModel += fmt.Sprintf(`"%s": "%s", `, key, value[0])
		}
		tempModel += `}`
		model := []byte(tempModel)
		ctx = ctx.WithModel(model)

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctrl.Response = w
		ctrl.Request = r

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := ctx.Logger.WithFields(logrus.Fields{
				DURATION:    duration,
				STATUS_CODE: statusCode,
				REMOTE:      ctx.RemoteAddress,
			})
			logger.Info(r.Method + SPACE + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				//TODO: Must render error page
				http.Error(w, "InternalServerError", 500)
			}
		}()

		w.Header().Set(CONTENT_TYPE, TEXT_HTML_CONTENT_TYPE)

		if err := f(ctx); err != nil {
			switch e := err.(type) {
			case app.ErrUnauthorized:
				// TODO: Must render 401 page
				http.Error(w, e.Error(), e.Code())
			case app.ErrForbidden:
				// TODO: Must render 403 page
				http.Error(w, e.Error(), e.Code())
			default:
				// TODO: Must render Error page
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}
	})
}

func (ctrl *Controller) HandleRest(f action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		defer r.Body.Close()

		ctx := ctrl.App.NewContext().
			WithRemoteAddress(ctrl.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger)

		model, _ := ioutil.ReadAll(r.Body)
		ctx = ctx.WithModel(model)

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctrl.Response = w
		ctrl.Request = r

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := ctx.Logger.WithFields(logrus.Fields{
				DURATION:    duration,
				STATUS_CODE: statusCode,
				REMOTE:      ctx.RemoteAddress,
			})
			logger.Info(r.Method + SPACE + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		w.Header().Set(CONTENT_TYPE, JSON_CONTENT_TYPE)

		httperror := func(w http.ResponseWriter, code int, message string) {
			w.Header().Set(CONTENT_TYPE, TEXT_PLAIN_CONTENT_TYPE)
			w.Header().Set(X_CONTENT_TYPE_OPTIONS, NO_SNIFF)
			w.WriteHeader(code)
			fmt.Fprint(w, message)
		}

		if err := f(ctx); err != nil {
			switch e := err.(type) {
			case app.ErrHTTP:
				httperror(w, e.Code(), e.Error())
			default:
				httperror(
					w,
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError),
				)
			}
		}
	})
}
