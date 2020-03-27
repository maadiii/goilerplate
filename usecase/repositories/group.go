package repositories

import (
	"goilerplate/domain/models"
)

type IGroupRepository interface {
	FindAll(*[]models.Group) error
	FindAlongRoles(*models.Group) error
	FindAlongRolesAndUsers(*models.Group) error
	Create(*models.Group) error
	Delete(*models.Group) error
	Update(*models.Group) error
}
