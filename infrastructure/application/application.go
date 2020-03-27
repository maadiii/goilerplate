package application

import (
	"goilerplate/domain/models"
	"goilerplate/infrastructure/datastore"

	"github.com/sirupsen/logrus"
)

var (
	Version = VERSION
	Name    = NAME
)

type Application struct {
	Config    *Config
	DBSession *datastore.Session
}

func (a *Application) NewContext() *Context {
	return &Context{
		Logger:    logrus.StandardLogger(),
		DBSession: a.DBSession,
	}
}

func New() (app *Application, err error) {
	app = &Application{}
	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := datastore.InitConfig()
	if err != nil {
		return nil, err
	}

	app.DBSession, err = datastore.NewSession(dbConfig)
	if err != nil {
		return nil, err
	}

	return app, err
}

func (a *Application) Close() error {
	return a.DBSession.Close()
}

func (a *Application) MigrateDB() {
	s := a.DBSession

	s.AutoMigrate(
		models.User{},
		models.Group{},
		models.Role{},
		models.GroupsRole{},
	)
	s.Model(&models.User{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	s.Model(&models.GroupsRole{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	s.Model(&models.GroupsRole{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
}

func (a *Application) InsertBaseData() {
	a.DBSession.Create(&user)
}

func (a *Application) DropDB() {
	a.DBSession.DropTableIfExists(models.GroupsRole{})
	a.DBSession.DropTableIfExists(models.User{}, models.Group{}, models.Role{})
}

var (
	user = models.User{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  []byte("123456"),
		Group: models.Group{
			Name:        "مدیر",
			Description: "به کل برنامه دسترسی دارد.",
			Roles: []models.Role{
				{
					EnName: "addGroup",
					FaName: "افزودن گروه کاربری",
				},
				{
					EnName: "editGroup",
					FaName: "ویرایش گروه کاربری",
				},
				{
					EnName: "deleteGroup",
					FaName: "حذف گروه کاربری",
				},
				{
					EnName: "groupsList",
					FaName: "مشاهده لیست گروه‌های کاربری",
				},
				{
					EnName: "groupView",
					FaName: "مشاهده گروه کاربری",
				},
				{
					EnName: "addUser",
					FaName: "افزودن کاربر",
				},
				{
					EnName: "editUser",
					FaName: "ویرایش کاربر",
				},
				{
					EnName: "deleteUser",
					FaName: "حذف کاربر",
				},
				{
					EnName: "usersList",
					FaName: "مشاهده لیست کاربران",
				},
				{
					EnName: "userView",
					FaName: "مشاهده کاربر",
				},
			},
		},
	}
)

const (
	EMPTY       = ""
	NAME        = "goilerplate"
	DESCRIPTION = "Goilerplate web application"
	VERSION     = "0.1.0dev"

	// config
	SECRET_KEY              = "secret-key"
	BLOCK_SECRET_KEY        = "block-secret-key"
	JWT_SECRET              = "jwt.secret"
	JWT_ALGORITHM           = "jwt.algorithm"
	JWT_MAXAGE              = "jwt.max_age"
	JWT_HTTPONLY            = "jwt.httponly"
	REFRESH_TOKEN_SECRET    = "jwt.refresh_token.secret"
	REFRESH_TOKEN_ALGORITHM = "jwt.refresh_token.algorithm"
	REFRESH_TOKEN_MAXAGE    = "jwt.refresh_token.max_age"
	REFRESH_TOKEN_SECURE    = "jwt.refresh_token.secure"
	REFRESH_TOKEN_HTTPONLY  = "jwt.refresh_token.httponly"
	REFRESH_TOKEN_PATH      = "jwt.refresh_token.path"
	STATIC                  = "static"
	CRT                     = "tls.crt"
	KEY                     = "tls.key"
	HS256                   = "HS256"
	HS384                   = "HS384"
	HS512                   = "HS512"

	// principles
	DIGITS   = "0123456789"
	SPECIALS = "~=+%^*/()[]{}/!@#$?|<>"
	ALL      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" +
		DIGITS + SPECIALS
	ALL_UNSPECIALS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" +
		DIGITS

	// Error messages
	UNAUTHORIZED = "Unauthorized"
	FORBIDDEN    = "Foribidden"
)
