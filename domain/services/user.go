package services

import (
	"goilerplate/app"
	"goilerplate/db"
	"goilerplate/domain/models"
)

type IUserService interface {
	Save(*models.User) error
	GetUserWithGroupAndRole(*models.User) error
}

type userService struct {
	dbsession *db.Session
}

func NewUserService(dbs *db.Session) IUserService {
	return &userService{
		dbsession: dbs,
	}
}

func (s *userService) Save(user *models.User) error {
	if s.dbsession.First(
		&models.User{},
		&models.User{MobileNumber: user.MobileNumber}).RecordNotFound() {
		return s.dbsession.Set(GORM_AUTOUPDATE, false).Create(user).Error
	} else {
		return app.NewErrConflict(USER_ALREADY_EXIST)
	}
}

func (s *userService) GetUserWithGroupAndRole(user *models.User) error {
	if s.dbsession.Preload("Group.Roles").
		Find(user, &models.User{ID: user.ID}).
		RecordNotFound() {
		return app.NewErrNotFound(USER_NOT_FOUND_WITH_ID)
	}

	return nil
}
