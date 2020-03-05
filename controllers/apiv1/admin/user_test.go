package admin_test

import (
	"goldfish/app/testcase"
	"goldfish/domain/models"
	"goldfish/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest("post", "/apiv1/admin/accounts")
	dbs := suittest.Context.DBSession

	group := models.Group{
		Name:        "TestGroup",
		Description: "TestGroup",
	}
	err := dbs.Create(&group).Error
	if err != nil {
		testcase.Fatal(err, t)
	}

	user := models.User{
		FirstName:    "AlreadyExit",
		MobileNumber: "09187710445",
		Password:     []byte("Mm1234"),
		GroupID:      group.ID,
	}
	err = dbs.Create(&user).Error
	if err != nil {
		testcase.Fatal(err, t)
	}

	t.Run("when phone number is blank", func(t *testing.T) {
		user := types.UserAdd{
			Password: "Mm123456",
			GroupID:  group.ID,
		}
		response := suittest.SendRestRequest(user)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"mobileNumber: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when phone number length is wrong", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "0",
			Password:     "Mm123456",
			GroupID:      group.ID,
		}
		response := suittest.SendRestRequest(user)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"mobileNumber: the length must be exactly 11.",
			response.Body.String(),
		)
	})

	t.Run("when phone number regex is wrong", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "aaaaaaaaaaa",
			Password:     "Mm123456",
			GroupID:      group.ID,
		}
		response := suittest.SendRestRequest(user)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"mobileNumber: must be in a valid format.",
			response.Body.String(),
		)
	})

	t.Run("when password is blank", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "09187710445",
			GroupID:      group.ID,
		}

		response := suittest.SendRestRequest(user)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"password: cannot be blank.",
			response.Body.String(),
		)
	})
	t.Run("when password is wrong format", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "09187710445",
			Password:     "1234",
			GroupID:      group.ID,
		}
		response := suittest.SendRestRequest(user)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"password: must be in a valid format.",
			response.Body.String(),
		)
	})

	t.Run("when user already exist", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "09187710445",
			Password:     "Mm1234",
			GroupID:      group.ID,
		}

		response := suittest.SendRestRequest(user)
		assert.Equal(t, 409, response.Code)
		assert.Equal(
			t,
			"User with this mobile number already exist",
			response.Body.String(),
		)
	})

	t.Run("when OK", func(t *testing.T) {
		user := types.UserAdd{
			MobileNumber: "09187710446",
			Password:     "Mm1234",
			GroupID:      group.ID,
		}
		response := suittest.SendRestRequest(user)
		assert.Equal(t, 200, response.Code)
	})
}
