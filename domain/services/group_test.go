package services_test

import (
	"goilerplate/app"
	"goilerplate/app/testcase"
	"goilerplate/domain/models"
	"goilerplate/domain/services"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession
	service := services.NewGroupService(suittest.Context.DBSession)

	t.Run("when group already exist", func(t *testing.T) {
		mockupGroup := models.Group{Name: "AlreadyExist"}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		err := service.Save(&models.Group{Name: "AlreadyExist"})
		assert.Equal(
			t,
			app.NewErrConflict("Group with this name already exist"),
			err,
		)
	})
	t.Run("when OK", func(t *testing.T) {
		err := service.Save(&models.Group{Name: "TestGroup"})
		assert.Equal(t, nil, err)
	})
}

func TestDeleteGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession
	service := services.NewGroupService(suittest.Context.DBSession)

	t.Run("when record not found", func(t *testing.T) {
		err := service.Delete(uuid.New())
		assert.Equal(
			t,
			app.NewErrNotFound("record not found"),
			err,
		)
	})

	t.Run("when group is admin", func(t *testing.T) {
		mockupGroup := models.Group{
			Name:        "Admin",
			Description: "Description",
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		err := service.Delete(mockupGroup.ID)
		assert.Equal(
			t,
			app.NewErrCustom(700, "You can't remove Admin group"),
			err,
		)
	})

	t.Run("when group has user", func(t *testing.T) {
		mockupUser := models.User{
			MobileNumber: "09187710445",
			Password:     []byte("123456"),
			Group: models.Group{
				Name:        "GroupName",
				Description: "Description",
			},
		}
		if err := dbsession.Create(&mockupUser).Error; err != nil {
			testcase.Fatal(err, t)
		}

		err := service.Delete(mockupUser.Group.ID)
		assert.Equal(
			t,
			app.NewErrCustom(701, "Group has users can't be removed"),
			err,
		)
	})

	t.Run("when OK", func(t *testing.T) {
		mockupGroup := models.Group{
			Name:        "GroupNameOK",
			Description: "Description",
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		err := service.Delete(mockupGroup.ID)
		assert.Equal(t, nil, err)
	})
}

func TestGroupList(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession
	for i := 0; i < 5; i++ {
		mockupGroup := models.Group{
			Name:        "TestName" + strconv.Itoa(i),
			Description: "Description" + strconv.Itoa(i),
		}

		err := dbsession.Create(&mockupGroup).Error
		if err != nil {
			testcase.Fatal(err, t)
		}
	}

	groups := []models.Group{}
	services.NewGroupService(dbsession).List(&groups)
	assert.Equal(t, 5, len(groups))
	assert.Equal(t, "TestName4", groups[4].Name)
}

func TestUpdateGroup(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession.
		Set("gorm:association_autoupdate", false)
	service := services.NewGroupService(suittest.Context.DBSession)

	mockupRole := models.Role{
		EnName: "EnName",
		FaName: "FaName",
	}
	if err := dbsession.Create(&mockupRole).Error; err != nil {
		testcase.Fatal(err, t)
	}

	t.Run("when group not found", func(t *testing.T) {
		err := service.Update(&models.Group{ID: uuid.New()})
		assert.Equal(t, app.NewErrNotFound("record not found"), err)
	})

	t.Run("when group is Admin", func(t *testing.T) {
		mockupGroup := models.Group{
			Name:        "Admin",
			Description: "Description",
			Roles:       []models.Role{{ID: mockupRole.ID}},
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		err := service.Update(&mockupGroup)
		assert.Equal(
			t,
			app.NewErrCustom(700, "You can't update Admin group"),
			err,
		)
	})

	t.Run("when group name is repetitious", func(t *testing.T) {
		mockupGroups := []*models.Group{
			&models.Group{
				Name:        "Repetitious",
				Description: "Description",
			},
			&models.Group{
				Name:        "GroupNameRep",
				Description: "Description",
			},
		}

		for _, group := range mockupGroups {
			if err := dbsession.Create(group).Error; err != nil {
				testcase.Fatal(err, t)
			}
		}

		group := models.Group{
			ID:   mockupGroups[1].ID,
			Name: "Repetitious",
		}
		err := service.Update(&group)
		assert.Equal(
			t,
			app.NewErrConflict("Group with this name already exist"),
			err,
		)
	})

	t.Run("when OK", func(t *testing.T) {
		mockupRole := models.Role{
			EnName: "afterUpdate",
			FaName: "faAfterUpdate",
		}
		if err := dbsession.Create(&mockupRole).Error; err != nil {
			testcase.Fatal(err, t)
		}

		mockupGroup := models.Group{
			Name:  "Updatable",
			Roles: []models.Role{{ID: mockupRole.ID}},
		}
		if err := dbsession.Create(&mockupGroup).Error; err != nil {
			testcase.Fatal(err, t)
		}

		updatedGroup := models.Group{
			ID:          mockupGroup.ID,
			Name:        "UpdatedGroup",
			Description: "DescriptionUpdated",
			Roles:       []models.Role{{ID: mockupRole.ID}},
		}
		err := service.Update(&updatedGroup)
		assert.Equal(t, nil, err)

		dbsession := suittest.Context.DBSession

		roles := []models.Role{}
		if err := dbsession.Find(&roles).Error; err != nil {
			testcase.Fatal(err, t)
		}
		assert.Equal(t, 2, len(roles))

		group := models.Group{}
		err = dbsession.Preload("Roles").
			First(&group, &models.Group{ID: updatedGroup.ID}).Error
		if err != nil {
			testcase.Fatal(err, t)
		}

		assert.Equal(t, "UpdatedGroup", group.Name, "name not equal")
		assert.Equal(
			t,
			"DescriptionUpdated",
			group.Description, "description not equal",
		)
		assert.Equal(t, 1, len(group.Roles), "len of roles not equal")
		assert.Equal(
			t,
			"afterUpdate",
			group.Roles[0].EnName,
			"role name not equal",
		)
		assert.Equal(
			t,
			"faAfterUpdate",
			group.Roles[0].FaName,
			"role name not equal",
		)
	})
}

func TestSelectGroupAndRoles(t *testing.T) {
	suittest.Init(t)
	dbsession := suittest.Context.DBSession

	mockupGroup := models.Group{
		Name:        "TestGroup",
		Description: "TestDescription",
		Roles: []models.Role{
			{
				EnName: "EnName",
				FaName: "FaName",
			},
		},
	}
	if err := dbsession.Create(&mockupGroup).Error; err != nil {
		testcase.Fatal(err, t)
	}

	service := services.NewGroupService(dbsession)

	t.Run("when group not found", func(t *testing.T) {
		err := service.SelectAndRoles(&models.Group{ID: uuid.New()})
		assert.Equal(
			t,
			app.NewErrNotFound("record not found"),
			err,
		)
	})

	t.Run("when OK", func(t *testing.T) {
		group := models.Group{ID: mockupGroup.ID}
		err := services.NewGroupService(dbsession).SelectAndRoles(&group)
		assert.NoError(t, err)
		assert.Equal(t, "TestGroup", group.Name)
		assert.Equal(t, "TestDescription", group.Description)
		assert.Len(t, group.Roles, 1)
		assert.Equal(t, "FaName", group.Roles[0].FaName)
	})
}
