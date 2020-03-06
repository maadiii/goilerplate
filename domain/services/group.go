package services

import (
	"goilerplate/app"
	"goilerplate/db"
	"goilerplate/domain/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type IGroupService interface {
	Save(*models.Group) error
	List(*[]models.Group)
	Delete(id uuid.UUID) error
	Update(*models.Group) error
	SelectAndRoles(*models.Group) error
}

type groupService struct {
	dbsession *db.Session
}

func NewGroupService(dbs *db.Session) IGroupService {
	return &groupService{
		dbsession: dbs,
	}
}

func (s *groupService) Save(g *models.Group) error {
	available := models.Group{Name: g.Name}
	err := s.dbsession.First(&available, available).Error
	if gorm.IsRecordNotFoundError(err) {
		return s.dbsession.Set(GORM_AUTOUPDATE, false).Create(g).Error
	}

	if available.ID != uuid.Nil {
		return app.NewErrConflict(GROUP_ALREADY_EXIST)
	}

	return err
}

func (s *groupService) List(groups *[]models.Group) {
	s.dbsession.Order("name").Find(groups)
}

func (s *groupService) Delete(id uuid.UUID) error {
	group := models.Group{}
	if !s.dbsession.
		First(&group, models.Group{ID: id}).
		Related(&group.Users, "Users").
		RecordNotFound() {
		if group.Name == ADMIN {
			return app.NewErrCustom(700, ADMIN_REMOVE_PERMISSION)
		}
		if len(group.Users) != 0 {
			return app.NewErrCustom(701, GROUP_HAS_USERS)
		}

		return s.dbsession.Delete(
			&models.GroupsRole{},
			&models.GroupsRole{GroupID: id},
		).Delete(
			&models.Group{},
			models.Group{ID: id},
		).Error
	} else {
		return app.NewErrNotFound(RECORD_NOT_FOUND)
	}
}

func (s *groupService) Update(g *models.Group) error {
	alreadyGroupExist := models.Group{}
	if !s.dbsession.
		First(&alreadyGroupExist, &models.Group{Name: g.Name}).
		RecordNotFound() {
		if alreadyGroupExist.Name == g.Name && alreadyGroupExist.ID != g.ID {
			return app.NewErrConflict(GROUP_ALREADY_EXIST)
		}
	}

	group := models.Group{}
	if !s.dbsession.First(&group, &models.Group{ID: g.ID}).RecordNotFound() {
		if group.Name == "Admin" {
			return app.NewErrCustom(700, ADMIN_UPDATE_PERMISSION)
		}
		s.dbsession.Unscoped().
			Delete(&models.GroupsRole{}, "group_id = ?", g.ID)

		return s.dbsession.Set(GORM_AUTOUPDATE, false).
			Model(g).Updates(&g).Error
	} else {
		return app.NewErrNotFound(RECORD_NOT_FOUND)
	}
}

func (s *groupService) SelectAndRoles(g *models.Group) error {
	err := s.dbsession.Preload("Roles").First(g, g).Error
	if gorm.IsRecordNotFoundError(err) {
		return app.NewErrNotFound(RECORD_NOT_FOUND)
	}
	return err
}
