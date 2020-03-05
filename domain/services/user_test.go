package services_test

import (
	"goldfish/app"
	"goldfish/app/testcase"
	"goldfish/domain/models"
	"goldfish/domain/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	users = []*models.User{}
)

func TestSaveUser(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession
	service := services.NewUserService(dbsession)

	group := models.Group{
		Name:        "TestGroup",
		Description: "TestDescription",
		Roles: []models.Role{
			{
				EnName: "EnTestRole",
				FaName: "FaTestRole",
			},
		},
	}
	if err := dbsession.Create(&group).Error; err != nil {
		testcase.Fatal(err, t)
	}

	user := &models.User{
		MobileNumber: "09187710445",
		Password:     []byte("Password"),
		GroupID:      group.ID,
	}
	err := dbsession.Create(&user).Error
	if err != nil {
		testcase.Fatal(err, t)
	}

	t.Run("when user already exist", func(t *testing.T) {
		err := service.Save(&models.User{MobileNumber: "09187710445"})
		assert.Equal(
			t,
			app.NewErrConflict("User with this mobile number already exist"),
			err,
		)
	})
	t.Run("when OK", func(t *testing.T) {
		err := service.Save(&models.User{
			MobileNumber: "09123456789",
			Password:     []byte("123456"),
			GroupID:      group.ID,
		})
		assert.NoError(t, err)
	})
}

func TestGetUserWithGroupAndRole(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession

	group := models.Group{
		Name:        "GroupName",
		Description: "Description",
		Roles: []models.Role{
			{
				EnName: "enRoleName",
				FaName: "faRoleName",
			},
		},
	}
	if err := dbsession.Create(&group).Error; err != nil {
		testcase.Fatal(err, t)
	}

	user := models.User{
		MobileNumber: "09187710445",
		Password:     []byte("123456"),
		GroupID:      group.ID,
	}
	if err := dbsession.Create(&user).Error; err != nil {
		testcase.Fatal(err, t)
	}

	service := services.NewUserService(dbsession)

	t.Run("when user not found", func(t *testing.T) {
		err := service.GetUserWithGroupAndRole(&models.User{ID: uuid.New()})
		assert.Equal(
			t,
			app.NewErrNotFound("User with this id not found"),
			err,
		)
	})
	t.Run("when OK", func(t *testing.T) {
		user := models.User{ID: user.ID}
		err := service.GetUserWithGroupAndRole(&user)
		assert.NoError(t, err)
		assert.Equal(t, "09187710445", user.MobileNumber)
		assert.Equal(t, user.Group.Roles[0].EnName, "enRoleName")
	})
}
