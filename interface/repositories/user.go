package repositories

import (
	"goilerplate/domain/models"
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/datastore"
	ur "goilerplate/usecase/repositories"
)

type userRepository struct {
	session *datastore.Session
}

func NewUserRepository(dbsession *datastore.Session) ur.IUserRepository {
	return &userRepository{dbsession}
}

func (r *userRepository) Create(u *models.User) error {
	if r.session.Where(MOBILE_NUMBER_CONDITION, u.MobileNumber).First(u).RecordNotFound() {
		return r.session.Set(GORM_AUTOUPDATE, false).Create(u).Error
	} else {
		return application.NewErrConflict(USER_ALREADY_EXIST)
	}
}

func (r *userRepository) FindAllAlongGroup(users *[]models.User, p int, s string) error {
	if s != "" {
		err := r.session.Where(FIRSTNAME_CONDITION, s).
			Or(LASTNAME_CONDITION, s).Or(MOBILE_NUMBER_CONDITION, s).Preload(GROUP).
			Limit(10).Offset(10 * (p - 1)).Find(&users).Error
		return err
	} else {
		err := r.session.Preload("Group").Limit(10).Offset(10 * (p - 1)).Find(&users).Error
		return err
	}
}

func (r *userRepository) Count(s string, c *int) error {
	if s == "" {
		return r.session.Table(USER_RELATION).Count(c).Error
	} else {
		return r.session.Model(&models.User{}).Where(FIRSTNAME_CONDITION, s).
			Or(LASTNAME_CONDITION, s).Or(MOBILE_NUMBER_CONDITION, s).Count(c).Error
	}
}

const (
	USER_ALREADY_EXIST      = "User with this mobile number already exist"
	FIRSTNAME_CONDITION     = "first_name = ?"
	LASTNAME_CONDITION      = "last_name = ?"
	MOBILE_NUMBER_CONDITION = "mobile_number = ?"
	USERS                   = "Users"
	USER_RELATION           = "users"
)
