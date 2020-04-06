package repositories

import "goilerplate/domain/models"

type IUserRepository interface {
	FindAllAlongGroup(*[]models.User, int, string) error
	Create(*models.User) error
	Count(string, *int) error
}
