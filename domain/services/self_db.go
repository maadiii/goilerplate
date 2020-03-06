package services

import (
	"goilerplate/db"
	"goilerplate/domain/models"
)

var (
	user = models.User{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  []byte("123456"),
		Group: models.Group{
			Name:        "Admin",
			Description: "دسترسی به کل سایت.",
			Roles: []models.Role{
				models.Role{
					EnName: "addUser",
					FaName: "افزودن کاربر",
				},
				models.Role{
					EnName: "addGroup",
					FaName: "افزودن گروه کاربری",
				},
				models.Role{
					EnName: "groupList",
					FaName: "لیست گروه‌های کاربری",
				},
				models.Role{
					EnName: "deleteGroup",
					FaName: "حذف گروه کاربری",
				},
			},
		},
	}
)

func InsertBaseData(dbs *db.Session) {
	dbs.Create(&user)
}

func DropDB(dbs *db.Session) {
	dbs.DropTableIfExists(models.GroupsRole{})
	dbs.DropTableIfExists(models.User{}, models.Group{}, models.Role{})
}

func MigrateDB(dbs *db.Session) {
	dbs.AutoMigrate(
		models.User{},
		models.Group{},
		models.Role{},
		models.GroupsRole{},
	)
	dbs.Model(&models.User{}).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	dbs.Model(&models.GroupsRole{}).
		AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	dbs.Model(&models.GroupsRole{}).
		AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
}
