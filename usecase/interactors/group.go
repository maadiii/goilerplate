package interactors

import (
	"goilerplate/domain/models"
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/utils"
	"goilerplate/usecase/presenters"
	"goilerplate/usecase/repositories"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type IGroupInteractor interface {
	All() ([]presenters.GroupPresent, error)
	GetAlongSelectedRoles(uuid.UUID, *[]presenters.RolePresent) (presenters.GroupAlongSelectedRolePresent, error)
	GetAlongRolesAndUsers(uuid.UUID) (presenters.GroupAlongRolesAndUsersPresent, error)
	Save(*GroupSave) (presenters.GroupPresent, error)
	Remove(*GroupDelete) error
	Edit(*GroupEdit) (presenters.GroupPresent, error)
}

type groupInteractor struct {
	repository repositories.IGroupRepository
	presenter  presenters.IGroupPresenter
}

func NewGroupInteractor(
	r repositories.IGroupRepository,
	p presenters.IGroupPresenter,
) IGroupInteractor {
	return &groupInteractor{r, p}
}

func (i *groupInteractor) All() ([]presenters.GroupPresent, error) {
	groups := []models.Group{}
	err := i.repository.FindAll(&groups)
	if err != nil {
		return nil, err
	}

	return i.presenter.PresentAll(&groups), nil
}

func (i *groupInteractor) GetAlongSelectedRoles(id uuid.UUID, roles *[]presenters.RolePresent,
) (presenters.GroupAlongSelectedRolePresent, error) {
	grps := models.Group{ID: id}
	if err := i.repository.FindAlongRoles(&grps); err != nil {
		return presenters.GroupAlongSelectedRolePresent{}, err
	}

	return i.presenter.PresentAlongSelectedRoles(&grps, roles), nil
}

func (i *groupInteractor) Save(g *GroupSave) (presenters.GroupPresent, error) {
	var groupPresent presenters.GroupPresent
	if err := g.Validate(); err != nil {
		return groupPresent, application.NewErrValidation(err.Error())
	}
	group := models.Group{
		Name:        g.Name,
		Description: g.Description,
		Roles:       make([]models.Role, len(g.RoleIds)),
	}

	for i, id := range g.RoleIds {
		group.Roles[i] = models.Role{ID: id}
	}

	err := i.repository.Create(&group)
	if err != nil {
		return groupPresent, err
	}

	groupPresent = i.presenter.PresentSave(&group)
	return groupPresent, nil
}

func (i *groupInteractor) Edit(g *GroupEdit) (presenters.GroupPresent, error) {
	groupPresent := presenters.GroupPresent{}
	if err := g.Validate(); err != nil {
		return groupPresent, application.NewErrValidation(err.Error())
	}
	group := models.Group{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Roles:       make([]models.Role, len(g.RoleIds)),
	}

	for i, id := range g.RoleIds {
		group.Roles[i] = models.Role{ID: id}
	}

	err := i.repository.Update(&group)
	if err != nil {
		return groupPresent, err
	}

	groupPresent = i.presenter.PresentEdit(&group)
	return groupPresent, nil
}

func (i *groupInteractor) Remove(g *GroupDelete) error {
	if err := utils.CheckUUID(g.ID); err != nil {
		return application.NewErrValidation(err.Error())
	}
	group := models.Group{ID: g.ID}
	return i.repository.Delete(&group)
}

func (i *groupInteractor) GetAlongRolesAndUsers(id uuid.UUID) (presenters.GroupAlongRolesAndUsersPresent, error) {
	group := presenters.GroupAlongRolesAndUsersPresent{}
	g := models.Group{ID: id}
	if err := i.repository.FindAlongRolesAndUsers(&g); err != nil {
		return group, err
	}

	return i.presenter.PresentAlongRolesAndUsers(&g), nil
}

type GroupSave struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	RoleIds     []uuid.UUID `json:"roleIds"`
}

type GroupEdit struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	RoleIds     []uuid.UUID `json:"roleIds"`
}

type GroupDelete struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (g GroupSave) Validate() error {
	return validation.ValidateStruct(
		&g,
		validation.Field(&g.Name, validation.Required, validation.Length(6, 128)),
		validation.Field(&g.Description, validation.Required, validation.Length(32, 256)),
		validation.Field(
			&g.RoleIds, validation.Required, validation.Length(1, 100),
			validation.Each(validation.By(utils.CheckUUID)),
		),
	)
}

func (g GroupEdit) Validate() error {
	return validation.ValidateStruct(
		&g,
		validation.Field(
			&g.ID, validation.Required, validation.Length(36, 36), validation.By(utils.CheckUUID),
		),
		validation.Field(&g.Name, validation.Required, validation.Length(6, 128)),
		validation.Field(&g.Description, validation.Required, validation.Length(32, 256)),
		validation.Field(&g.RoleIds, validation.Required, validation.Length(1, 100),
			validation.Each(validation.By(utils.CheckUUID)),
		),
	)
}
