package repositories

import "goilerplate/domain/models"

type IRoleRepository interface {
	FindAll(*[]models.Role) error
}
