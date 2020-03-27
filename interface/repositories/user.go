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
	if r.session.Where(MOBILE_NUMBER, u.MobileNumber).First(u).RecordNotFound() {
		return r.session.Set(GORM_AUTOUPDATE, false).Create(u).Error
	} else {
		return application.NewErrConflict(USER_ALREADY_EXIST)
	}
}

const (
	MOBILE_NUMBER      = "mobile_number = ?"
	USER_ALREADY_EXIST = "User with this mobile number already exist"
)
