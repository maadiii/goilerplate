package controller

import (
	"fmt"
	"goilerplate/infrastructure/application"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

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
		return config, fmt.Errorf("%s %s", PRE_ERROR, HTTP_ERROR)
	}
	if config.HTTPSPort == 0 {
		return config, fmt.Errorf("%s %s", PRE_ERROR, HTTPS_ERROR)
	}

	return config, nil
}

type Controller struct {
	Application *application.Application
	Config      *Config
}

type RestController struct {
	*Controller
}

func NewController(a *application.Application) (controller *Controller, err error) {
	controller = &Controller{Application: a}
	controller.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func NewRestController(controller *Controller) *RestController {
	return &RestController{controller}
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

type action func(*application.Context) error

//func (ctrl *Controller) Authorize(f action, role string) action {
//	return func(ctx *application.Context) error {
//		tokenString, err := ctrl.Application.ReadCookie(ctrl.Request, ACCESS_COOKIE)
//		if err != nil {
//			return application.NewErrUnauthorized()
//		}
//
//		claims := &application.Claims{}
//		keyFunc := func(token *jwt.Token) (interface{}, error) {
//			return ctrl.Application.Config.JWT.Secret, nil
//		}
//		token, err := jwt.ParseWithClaims(
//			tokenString,
//			claims,
//			keyFunc,
//		)
//		if err != nil {
//			return application.NewErrUnauthorized()
//		}
//
//		if !token.Valid {
//			refreshTokenString, err := ctrl.Application.
//				ReadCookie(ctrl.Request, REFRESH_COOKIE)
//			if err != nil {
//				return application.NewErrUnauthorized()
//			}
//
//			refreshClaims := &application.Claims{}
//			refreshToken, err := jwt.ParseWithClaims(
//				refreshTokenString,
//				refreshClaims,
//				keyFunc,
//			)
//			if err != nil {
//				return application.NewErrUnauthorized()
//			}
//
//			if !refreshToken.Valid {
//				return application.NewErrUnauthorized()
//			}
//
//			if refreshClaims.ID != claims.ID {
//				return application.NewErrUnauthorized()
//			}
//
//			user := models.User{ID: refreshClaims.ID}
//			err = services.NewUserService(ctx.DBSession).
//				GetUserWithGroupAndRole(&user)
//			if err != nil {
//				return application.NewErrUnauthorized()
//			}
//
//			roles := make([]string, len(user.Group.Roles))
//			for i, role := range user.Group.Roles {
//				roles[i] = role.EnName
//			}
//
//			jwt, err := ctrl.Application.CreateJWT(
//				user.ID,
//				user.FirstName,
//				user.LastName,
//				false,
//				roles...,
//			)
//			if err != nil {
//				return err
//			}
//
//			if err := ctrl.Application.SetCookie(
//				ACCESS_COOKIE,
//				jwt,
//				false,
//				ctrl.Response,
//			); err != nil {
//				return err
//			}
//
//			claims.Roles = roles
//		}
//
//		user := &models.User{
//			ID:        claims.ID,
//			FirstName: claims.FirstName,
//			LastName:  claims.LastName,
//		}
//		ctx = ctx.WithUser(user)
//
//		doNext := false
//		if role != EMPTY {
//			for _, r := range claims.Roles {
//				if r == role {
//					doNext = true
//					break
//				}
//			}
//			if !doNext {
//				return application.NewErrForbidden()
//			}
//		}
//
//		return f(ctx)
//	}
//}

func (ctrl *Controller) Handle(f action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		r.ParseForm()

		ctx := ctrl.Application.NewContext().WithRemoteAddress(ctrl.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger).WithRequest(r).WithResponseWriter(w)

		// TODO: handle model in form request
		//tempModel := `{`
		//for key, value := range r.Form {
		//	tempModel += fmt.Sprintf(`"%s": "%s", `, key, value[0])
		//}
		//tempModel += `}`
		//model := []byte(tempModel)
		//ctx = ctx.WithModel(model)

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				//TODO: Must render error page
				http.Error(w, "InternalServerError", 500)
			}
		}()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

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

		w.Header().Set(CONTENT_TYPE, TEXT_HTML_CONTENT_TYPE)

		if err := f(ctx); err != nil {
			switch e := err.(type) {
			case application.ErrUnauthorized:
				// TODO: Must render 401 page
				http.Error(w, e.Error(), e.Code())
			case application.ErrForbidden:
				// TODO: Must render 403 page
				http.Error(w, e.Error(), e.Code())
			case application.ErrNotFound:
				// TODO: Must render 404 page
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

func (ctrl *RestController) Handle(f action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)
		defer r.Body.Close()

		ctx := ctrl.Application.NewContext().WithRemoteAddress(ctrl.IPAddressForRequest(r))
		ctx = ctx.WithLogger(ctx.Logger).WithRequest(r).WithResponseWriter(w)

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

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
			case application.ErrHTTP:
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

const (
	EMPTY = ""
	SPACE = " "
	COMMA = ","
	COLON = ":"

	// methods
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	PATCH  = "PATCH"

	// root
	HTTPS_PORT              = "ports.https"
	HTTP_PORT               = "ports.http"
	DOMAIN_NAME             = "domain-name"
	PROXY_COUNT             = "proxy-count"
	DURATION                = "duration"
	STATUS_CODE             = "status_code"
	REMOTE                  = "remote"
	CONTENT_TYPE            = "Content-Type"
	JSON_CONTENT_TYPE       = "application/json"
	TEXT_PLAIN_CONTENT_TYPE = "text/plain; charset=utf-8"
	TEXT_HTML_CONTENT_TYPE  = "text/html"
	NO_SNIFF                = "nosniff"
	X_CONTENT_TYPE_OPTIONS  = "X-Content-Type-Options"
	X_FORWARDED_FOR         = "X-Forwarded-For"
	ACCESS_COOKIE           = "access"
	REFRESH_COOKIE          = "refresh"

	// error messages
	HTTP_ERROR  = "ports.http is not set in config file"
	HTTPS_ERROR = "ports.https is not set in config file"
)
