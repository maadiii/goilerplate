package testcase

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goldfish/app"
	"goldfish/controllers"
	"goldfish/controllers/route"
	"goldfish/db"
	"goldfish/domain/services"

	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type T struct {
	Controller  *controllers.Controller
	Context     *app.Context
	Testing     *testing.T
	httpRequest *http.Request
	httpRoute   *httprouter.Router
}

func removeTestDB(config *db.Config) error {
	errmsg := `app/testcase.go removeTestDB(), %v`

	var dbname string
	splited := strings.Split(config.Test, " ")
	for _, s := range splited {
		tmp := strings.Split(s, "=")
		if tmp[0] == "dbname" {
			dbname = tmp[1]
		}
	}

	dbs, err := sql.Open(config.Name, config.Admin)
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}
	defer dbs.Close()

	err = dbs.Ping()
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}

	statement := fmt.Sprintf(DROP_DATABASE_STATEMENT, dbname)
	_, err = dbs.Exec(statement)
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}

	return err
}

func createTestDB(config *db.Config) error {
	errmsg := `app/testcase.go createTestDB(), %v`

	var dbname string
	var dbuser string
	splited := strings.Split(config.Test, " ")
	for _, s := range splited {
		tmp := strings.Split(s, "=")
		if tmp[0] == "dbname" {
			dbname = tmp[1]
		} else if tmp[0] == "user" {
			dbuser = tmp[1]
		}
	}

	dbs, err := sql.Open(config.Name, config.Admin)
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}
	defer dbs.Close()

	err = dbs.Ping()
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}

	removeTestDB(config)

	statement := fmt.Sprintf(CREATE_DATABASE_STATEMENT, dbname)
	_, err = dbs.Exec(statement)
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}

	statement = fmt.Sprintf(SET_GRANT, dbname, dbuser)
	_, err = dbs.Exec(statement)
	if err != nil {
		return fmt.Errorf(errmsg, err)
	}

	return err
}

func testDBSession(config *db.Config) (*db.Session, error) {
	errmsg := `app/testcase.go NewTestDB(), %v`
	gdbs, err := gorm.Open(config.Name, config.Test)
	if err != nil {
		return nil, fmt.Errorf(errmsg, err)
	}

	return &db.Session{gdbs}, nil
}

type Mockup func(*db.Session)

func New() *T {
	appConfig, err := app.InitConfig()
	if err != nil {
		panic(err)
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		panic(err)
	}

	err = createTestDB(dbConfig)
	if err != nil {
		panic(err)
	}

	dbs, err := testDBSession(dbConfig)
	if err != nil {
		removeTestDB(dbConfig)
		panic(err)
	}

	app := &app.App{appConfig, dbs}
	c, err := controllers.New(app)
	if err != nil {
		removeTestDB(dbConfig)
		panic(err)
	}

	testctrl := &T{
		Controller: c,
		Context:    app.NewContext(),
	}

	return testctrl
}

func (t *T) Init(te *testing.T) {
	services.DropDB(t.Context.DBSession)
	services.MigrateDB(t.Context.DBSession)
	t.Testing = te
}

func (t *T) Close() {
	err := t.Context.DBSession.Close()
	if err != nil {
		panic(err)
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		panic(err)
	}

	if err := removeTestDB(dbConfig); err != nil {
		panic(err)
	}
}

func (t *T) authentication(r *http.Request) (*http.Request, error) {
	app := t.Controller.App
	jwt, err := app.CreateJWT(uuid.New(), EMPTY, EMPTY, false, EMPTY)
	if err != nil {
		return r, err
	}

	secureCookie := securecookie.
		New(app.Config.SecretKey, app.Config.BlockSecretKey)
	encoded, err := secureCookie.Encode(COOKIE, jwt)
	if err == nil {
		cookie := &http.Cookie{
			Name:  COOKIE,
			Value: encoded,
		}

		r.AddCookie(cookie)
	}

	return r, err
}

func (t *T) SendRestRequest(data interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	if data != nil {
		switch ty := data.(type) {
		case string:
			ty, _ = data.(string)
			t.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(ty)))
		default:
			d, _ := json.Marshal(&data)
			t.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(d))
		}
	}

	t.httpRoute.ServeHTTP(rr, t.httpRequest)

	return rr
}

func (t *T) SendAuthentRestRequest(
	data interface{},
	request *http.Request,
	router *httprouter.Router,
) *httptest.ResponseRecorder {
	request, err := t.authentication(request)
	if err != nil {
		panic(err)
	}

	return t.SendRestRequest(data)
}

func (t *T) SendBadJson(req *http.Request, router *httprouter.Router) {
	t.Testing.Run(SEND_BAD_JSON, func(te *testing.T) {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(BAD_JSON)))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(te, 400, rr.Code)
	})
}

func (t *T) SetHTTPRequest(method, url string) {
	var err error
	t.httpRequest, err = http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		panic(err)
	}

	t.httpRoute = httprouter.New()
	route.Route(t.Controller, t.httpRoute)
}

func Fatal(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
