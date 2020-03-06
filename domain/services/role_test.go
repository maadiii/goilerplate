package services_test

import (
	"goilerplate/domain/models"
	"goilerplate/domain/services"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func roleMockup() []*models.Role {
	dbs := suittest.Context.DBSession
	roles := []*models.Role{}
	for i := 0; i < 5; i++ {
		role := models.Role{
			EnName: "EnName" + strconv.Itoa(i),
			FaName: "FaName" + strconv.Itoa(i),
		}
		err := dbs.Create(&role).Error
		if err != nil {
			panic(err)
		}

		roles = append(roles, &role)
	}

	return roles
}

func TestGetRole(t *testing.T) {
	suittest.Init(t)
	roles := roleMockup()

	role := models.Role{ID: roles[0].ID}
	err := services.NewRoleService(suittest.Context.DBSession).Get(&role)

	assert.Equal(t, nil, err)
	assert.Equal(t, "EnName0", role.EnName)
}

func TestRoleList(t *testing.T) {
	suittest.Init(t)
	_ = roleMockup()
	list := []models.Role{}
	err := services.NewRoleService(suittest.Context.DBSession).List(&list)

	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(list))
}
