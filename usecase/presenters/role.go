package presenters

import (
	"goilerplate/domain/models"

	"github.com/google/uuid"
)

type IRolePresenter interface {
	PresentAll(*[]models.Role) []RolePresent
}

type RolePresent struct {
	ID     uuid.UUID `json:"id"`
	FaName string    `json:"faName"`
	EnName string    `json:"enName"`
}

type SelectedRolePresent struct {
	ID       uuid.UUID `json:"id"`
	FaName   string    `json:"faName"`
	EnName   string    `json:"enName"`
	Selected bool      `json:"selected"`
}
