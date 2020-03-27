package repositories

import "goilerplate/domain/models"

type IUserRepository interface {
	Create(*models.User) error
}
