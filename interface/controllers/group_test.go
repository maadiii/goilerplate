package controllers_test

import (
	"fmt"
	"goilerplate/domain/models"
	"goilerplate/infrastructure/testutil"
	"goilerplate/usecase/interactors"
	"net/http"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/google/uuid"
)

func TestAddGroup(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest(http.MethodGet, "/groups/create")
	res := suittest.SendViewRequest(nil)
	assert.Equal(t, 200, res.Code)
}

func TestViewGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Controller.GetBase().Application.DBSession
	group := &models.Group{
		Name:        "TestName",
		Description: "Description",
	}
	dbsession.Create(&group)

	t.Run("when id is invalid", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/view/bad_id")
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 500, res.Code)
	})

	// TODO: fix it for not found error page
	t.Run("when group not found", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/view/"+uuid.New().String())
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 404, res.Code)
	})

	t.Run("when OK", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/view/"+group.ID.String())
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 200, res.Code)
	})
}

func TestEditGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Controller.GetBase().Application.DBSession
	group := &models.Group{
		Name:        "TestName",
		Description: "Description",
		Roles: []models.Role{
			{
				FaName: "faName",
				EnName: "enName",
			},
		},
	}
	dbsession.Create(&group)

	// TODO: fix it for not found error page
	t.Run("when group not found", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/edit/"+uuid.New().String())
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 404, res.Code)
	})

	t.Run("when id is invalid", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/edit/bad_id")
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 500, res.Code)
	})

	t.Run("when OK", func(t *testing.T) {
		suittest.SetHTTPRequest(http.MethodGet, "/groups/edit/"+group.ID.String())
		res := suittest.SendViewRequest(nil)
		assert.Equal(t, 200, res.Code)
	})
}

func TestGroupList(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest(http.MethodGet, "/groups")
	res := suittest.SendViewRequest(nil)
	assert.Equal(t, 200, res.Code)
}

func TestPostGroup(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest(http.MethodPost, "/apiv1/groups")
	dbsession := suittest.Controller.GetBase().Application.DBSession

	role := models.Role{
		FaName: "FaName",
		EnName: "EnName",
	}
	err := dbsession.Create(&role).Error
	if err != nil {
		testutil.Fatal(err, t)
	}

	err = dbsession.Create(&models.Group{Name: "AlreadyExist"}).Error
	if err != nil {
		testutil.Fatal(err, t)
	}

	t.Run("when name is blank", func(t *testing.T) {
		group := interactors.GroupSave{
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"name: cannot be blank.",
			response.Body.String(),
		)
	})
	t.Run("when name length is wrong", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        "a",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"name: the length must be between 6 and 128.",
			response.Body.String(),
		)
	})

	t.Run("when description is blank", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:    "testname",
			RoleIds: []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"description: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when description length is wrong", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        "testname2",
			Description: strings.Repeat("#", 31),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"description: the length must be between 32 and 256.",
			response.Body.String(),
		)
	})

	t.Run("when length of roles is 0", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        strings.Repeat("#", 6),
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roleIds: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when each of roles is invalid", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        strings.Repeat("#", 6),
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.UUID{}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roleIds: (0: invalid uuid.).",
			response.Body.String(),
		)
	})

	t.Run("when group already exist", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        "AlreadyExist",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 409, response.Code)
		assert.Equal(
			t,
			"Group with this name already exist",
			response.Body.String(),
		)
	})

	t.Run("when OK", func(t *testing.T) {
		group := interactors.GroupSave{
			Name:        "TestName1",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{role.ID},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 200, response.Code)
	})
}

func TestPutGroup(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest(http.MethodPut, "/apiv1/groups")
	dbsession := suittest.Controller.GetBase().Application.DBSession

	t.Run("when name is blank", func(t *testing.T) {
		group := interactors.GroupEdit{
			ID:          uuid.New(),
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"name: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when name length is wrong", func(t *testing.T) {
		group := interactors.GroupEdit{
			ID:          uuid.New(),
			Name:        "a",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"name: the length must be between 6 and 128.",
			response.Body.String(),
		)
	})

	t.Run("when description is blank", func(t *testing.T) {
		group := interactors.GroupEdit{
			ID:      uuid.New(),
			Name:    "GroupName",
			RoleIds: []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"description: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when roles is blank", func(t *testing.T) {
		group := interactors.GroupEdit{
			ID:          uuid.New(),
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roleIds: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when one of roles is invalid", func(t *testing.T) {
		group := interactors.GroupEdit{
			ID:          uuid.New(),
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.UUID{}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roleIds: (0: invalid uuid.).",
			response.Body.String(),
		)
	})

	t.Run("when id is nil", func(t *testing.T) {
		group := interactors.GroupEdit{
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"id: invalid uuid.",
			response.Body.String(),
		)
	})

	t.Run("when group is repetitous", func(t *testing.T) {
		mockupGroups := []*models.Group{
			{Name: "Repetitous"},
			{Name: "Patchable"},
		}
		for _, role := range mockupGroups {
			if err := dbsession.Create(role).Error; err != nil {
				testutil.Fatal(err, t)
			}
		}

		group := interactors.GroupEdit{
			ID:          mockupGroups[1].ID,
			Name:        mockupGroups[0].Name,
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}
		response := suittest.SendRestRequest(group)
		assert.Equal(t, 409, response.Code)
		assert.Equal(
			t,
			"Group with this name already exist",
			response.Body.String(),
		)
	})

	t.Run("when group is admin", func(t *testing.T) {
		mockupGroup := models.Group{Name: "مدیر"}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testutil.Fatal(err, t)
		}

		group := interactors.GroupEdit{
			ID:          mockupGroup.ID,
			Name:        "NewName",
			Description: strings.Repeat("#", 32),
			RoleIds:     []uuid.UUID{uuid.New()},
		}
		response := suittest.SendRestRequest(group)
		assert.Equal(t, 700, response.Code)
		assert.Equal(
			t,
			"You can't update Admin group",
			response.Body.String(),
		)
	})

	t.Run("when OK", func(t *testing.T) {
		mockupRoles := []*models.Role{
			{
				EnName: "enName0",
				FaName: "faName0",
			},
			{
				EnName: "enName1",
				FaName: "faName1",
			},
		}
		for _, role := range mockupRoles {
			if err := dbsession.Create(role).Error; err != nil {
				testutil.Fatal(err, t)
			}
		}

		mockupGroup := models.Group{
			Name:        "GroupName",
			Description: "Description",
			Roles:       []models.Role{{ID: mockupRoles[0].ID}},
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testutil.Fatal(err, t)
		}

		group := interactors.GroupEdit{
			ID:          mockupGroup.ID,
			Name:        "UpdatedGroupName",
			Description: strings.Repeat("UpdatedDescription", 3),
			RoleIds:     []uuid.UUID{mockupRoles[1].ID},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 200, response.Code)

		updatedGroup := models.Group{ID: mockupGroup.ID}
		if err := dbsession.Preload("Roles").
			First(&updatedGroup, updatedGroup).Error; err != nil {
			testutil.Fatal(err, t)
		}
		assert.Equal(t, "UpdatedGroupName", updatedGroup.Name)
		assert.Equal(
			t,
			strings.Repeat("UpdatedDescription", 3),
			updatedGroup.Description,
		)
		assert.Equal(t, 1, len(updatedGroup.Roles))
		assert.Equal(t, "enName1", updatedGroup.Roles[0].EnName)
	})
}

func TestDeleteGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Controller.GetBase().Application.DBSession
	suittest.SetHTTPRequest(http.MethodDelete, "/apiv1/groups")

	t.Run("when group not found", func(t *testing.T) {
		id := fmt.Sprintf(`{"id":"%v"}`, uuid.New())
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 404, response.Code)
		assert.Equal(
			t,
			"record not found",
			response.Body.String(),
		)
	})

	t.Run("when id is invalid format", func(t *testing.T) {
		id := fmt.Sprintf(`{"id": "%v"}`, uuid.Nil)
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"invalid uuid",
			response.Body.String(),
		)
	})

	t.Run("when group is admin", func(t *testing.T) {
		group := models.Group{Name: "مدیر"}
		dbsession.Create(&group)
		id := fmt.Sprintf(`{"id":"%v"}`, group.ID)
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 700, response.Code)
		assert.Equal(
			t,
			"You can't remove Admin group",
			response.Body.String(),
		)
	})

	t.Run("when has some users", func(t *testing.T) {
		user := models.User{
			MobileNumber: "09187710445",
			Password:     []byte("123456"),
			Group:        models.Group{Name: "HasUser"},
		}
		dbsession.Create(&user)
		id := fmt.Sprintf(`{"id":"%v"}`, user.Group.ID)
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 701, response.Code)
		assert.Equal(
			t,
			"Group has users can't be removed",
			response.Body.String(),
		)
	})

	t.Run("when OK", func(t *testing.T) {
		group := models.Group{Name: "TestGroup"}
		dbsession.Create(&group)
		id := fmt.Sprintf(`{"id":"%v"}`, group.ID)
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 200, response.Code)
	})
}
