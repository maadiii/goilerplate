package presenters

import (
	"goilerplate/domain/models"

	"github.com/google/uuid"
)

type IGroupPresenter interface {
	PresentAll(*[]models.Group) []GroupPresent
	PresentAlongSelectedRoles(*models.Group, *[]RolePresent) GroupAlongSelectedRolePresent
	PresentAlongRolesAndUsers(*models.Group) GroupAlongRolesAndUsersPresent
	PresentSave(*models.Group) GroupPresent
	PresentEdit(*models.Group) GroupPresent
}

type GroupPresent struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type GroupAlongSelectedRolePresent struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Roles       []SelectedRolePresent `json:"roles"`
}

type GroupAlongRolesAndUsersPresent struct {
	Name        string
	Description string
	Roles       []RolePresent
	Users       []UserPresent
}
