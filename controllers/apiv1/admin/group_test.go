package admin_test

import (
	"fmt"
	"goilerplate/app/testcase"
	"goilerplate/domain/models"
	"goilerplate/types"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostGroup(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest("post", "/apiv1/admin/groups")
	dbs := suittest.Context.DBSession

	role := models.Role{
		FaName: "FaName",
		EnName: "EnName",
	}
	err := dbs.Create(&role).Error
	if err != nil {
		testcase.Fatal(err, t)
	}

	err = dbs.Create(&models.Group{Name: "AlreadyExist"}).Error
	if err != nil {
		testcase.Fatal(err, t)
	}

	t.Run("when name is blank", func(t *testing.T) {
		group := types.GroupAdd{
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		group := types.GroupAdd{
			Name:        "a",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		group := types.GroupAdd{
			Name:  "testname",
			Roles: []types.Role{{ID: uuid.New()}},
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
		group := types.GroupAdd{
			Name:        "testname2",
			Description: strings.Repeat("#", 31),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		group := types.GroupAdd{
			Name:        strings.Repeat("#", 6),
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roles: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when each of roles is invalid", func(t *testing.T) {
		group := types.GroupAdd{
			Name:        strings.Repeat("#", 6),
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roles: (0: invalid role.).",
			response.Body.String(),
		)
	})

	t.Run("when group already exist", func(t *testing.T) {
		group := types.GroupAdd{
			Name:        "AlreadyExist",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: role.ID}},
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
		group := types.GroupAdd{
			Name:        "TestName1",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: role.ID}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 200, response.Code)
	})
}

func TestDeleteGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession
	suittest.SetHTTPRequest("delete", "/apiv1/admin/groups")

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
			"id: cannot be nil.",
			response.Body.String(),
		)
	})

	t.Run("when group is admin", func(t *testing.T) {
		group := models.Group{Name: "Admin"}
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

	t.Run("when id is bad uuid", func(t *testing.T) {
		id := `{"id":"bad uuid"}`
		response := suittest.SendRestRequest(id)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"invalid UUID length: 8",
			response.Body.String(),
		)
	})
}

func TestPatch(t *testing.T) {
	suittest.Init(t)
	suittest.SetHTTPRequest("patch", "/apiv1/admin/groups")
	dbsession := suittest.Context.DBSession

	t.Run("when name is blank", func(t *testing.T) {
		group := types.GroupEdit{
			ID:          uuid.New(),
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		group := types.GroupEdit{
			ID:          uuid.New(),
			Name:        "a",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		group := types.GroupEdit{
			ID:    uuid.New(),
			Name:  "GroupName",
			Roles: []types.Role{{ID: uuid.New()}},
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
		group := types.GroupEdit{
			ID:          uuid.New(),
			Name:        "GroupName",
			Description: "a",
			Roles:       []types.Role{{ID: uuid.New()}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"description: the length must be between 32 and 256.",
			response.Body.String(),
		)
	})

	t.Run("when roles is blank", func(t *testing.T) {
		group := types.GroupEdit{
			ID:          uuid.New(),
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roles: cannot be blank.",
			response.Body.String(),
		)
	})

	t.Run("when one of roles is invalid", func(t *testing.T) {
		group := types.GroupEdit{
			ID:          uuid.New(),
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 400, response.Code)
		assert.Equal(
			t,
			"roles: (0: invalid role.).",
			response.Body.String(),
		)
	})

	t.Run("when id is nil", func(t *testing.T) {
		group := types.GroupEdit{
			Name:        "GroupName",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
				testcase.Fatal(err, t)
			}
		}

		group := types.GroupEdit{
			ID:          mockupGroups[1].ID,
			Name:        mockupGroups[0].Name,
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
		mockupGroup := models.Group{Name: "Admin"}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		group := types.GroupEdit{
			ID:          mockupGroup.ID,
			Name:        "NewName",
			Description: strings.Repeat("#", 32),
			Roles:       []types.Role{{ID: uuid.New()}},
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
				testcase.Fatal(err, t)
			}
		}

		mockupGroup := models.Group{
			Name:        "GroupName",
			Description: "Description",
			Roles:       []models.Role{{ID: mockupRoles[0].ID}},
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		group := types.GroupEdit{
			ID:          mockupGroup.ID,
			Name:        "UpdatedGroupName",
			Description: strings.Repeat("UpdatedDescription", 3),
			Roles:       []types.Role{{ID: mockupRoles[1].ID}},
		}

		response := suittest.SendRestRequest(group)
		assert.Equal(t, 200, response.Code)

		updatedGroup := models.Group{ID: mockupGroup.ID}
		if err := dbsession.Preload("Roles").
			First(&updatedGroup, updatedGroup).Error; err != nil {
			testcase.Fatal(err, t)
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
