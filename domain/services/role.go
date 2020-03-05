package services

import (
	"goldfish/db"
	"goldfish/domain/models"
)

type IRoleService interface {
	List(*[]models.Role) error
	Get(*models.Role) error
}

type roleService struct {
	dbsession *db.Session
}

func NewRoleService(dbs *db.Session) IRoleService {
	return roleService{dbsession: dbs}
}

func (s roleService) List(l *[]models.Role) error {
	return s.dbsession.Order("created_at").Find(l).Error
}

func (s roleService) Get(r *models.Role) error {
	return s.dbsession.Find(r, &models.Role{ID: r.ID}).Error
}
