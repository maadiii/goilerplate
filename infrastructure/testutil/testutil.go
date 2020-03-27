package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"goilerplate/infrastructure/datastore"
	"goilerplate/interface/router"
	"goilerplate/registry"
	"goilerplate/usecase/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
)

type T struct {
	Controller  controllers.IRootController
	Testing     *testing.T
	httpRequest *http.Request
	httpRoute   *httprouter.Router
}

func removeTestDB(config *datastore.Config) error {
	errmsg := `infrastructure/testutil/testutil.go removeTestDB(), %v`

	var dbname string
	splited := strings.Split(config.Test, " ")
	for _, s := range splited {
		tmp := strings.Split(s, "=")
		if tmp[0] == "dbname" {
			dbname = tmp[1]
		}
	}

	dbs, err := sql.Open(config.Driver, config.Admin)
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

func createTestDB(config *datastore.Config) error {
	errmsg := `infrastructure/testutil/testutil.go createTestDB(), %v`

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

	dbs, err := sql.Open(config.Driver, config.Admin)
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

func New() *T {
	dbConfig, err := datastore.InitConfig()
	if err != nil {
		panic(err)
	}

	err = createTestDB(dbConfig)
	if err != nil {
		panic(err)
	}

	reg, err := registry.NewTestRegistry()
	if err != nil {
		panic(err)
	}
	c := reg.NewRootController()

	testctrl := &T{
		Controller: c,
	}

	return testctrl
}

func (t *T) Init(te *testing.T) {
	t.Controller.GetBase().Application.DropDB()
	t.Controller.GetBase().Application.MigrateDB()
	//t.Controller.GetBase().Application.InsertBaseData()
	t.Testing = te
}

func (t *T) Close() {
	err := t.Controller.GetBase().Application.Close()
	if err != nil {
		panic(err)
	}

	dbConfig, err := datastore.InitConfig()
	if err != nil {
		panic(err)
	}

	if err := removeTestDB(dbConfig); err != nil {
		panic(err)
	}
}

func (t *T) authentication(r *http.Request) (*http.Request, error) {
	app := t.Controller.GetBase().Application
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

func (t *T) SendViewRequest(data interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	// TODO: Set form data
	if data != nil {
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
	router.Route(t.httpRoute, t.Controller)
}

func Fatal(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

const (
	EMPTY = ""

	// testcase
	DROP_DATABASE_STATEMENT   = `DROP DATABASE %v;`
	CREATE_DATABASE_STATEMENT = `CREATE DATABASE %v;`
	SET_GRANT                 = `grant ALL privileges on database "%v" to "%v"`
	COOKIE                    = "access"
	SEND_BAD_JSON             = "when send bad json"
	BAD_JSON                  = `{"badjson}`
)
