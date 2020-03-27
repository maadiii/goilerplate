package repositories

import (
	"goilerplate/domain/models"
	"goilerplate/infrastructure/datastore"
	ur "goilerplate/usecase/repositories"
)

type roleRepository struct {
	session *datastore.Session
}

func NewRoleRepository(dbsession *datastore.Session) ur.IRoleRepository {
	return &roleRepository{session: dbsession}
}

func (repository *roleRepository) FindAll(roles *[]models.Role) error {
	return repository.session.Find(roles).Error
}
