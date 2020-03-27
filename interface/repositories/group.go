package repositories

import (
	"goilerplate/domain/models"
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/datastore"
	ur "goilerplate/usecase/repositories"

	"github.com/jinzhu/gorm"
)

type groupRepository struct {
	session *datastore.Session
}

func NewGroupRepository(dbsession *datastore.Session) ur.IGroupRepository {
	return &groupRepository{session: dbsession}
}

func (r *groupRepository) FindAll(groups *[]models.Group) error {
	return r.session.Order("name").Find(groups).Error
}

func (r *groupRepository) Find(g *models.Group) error {
	err := r.session.First(g, g).Error
	if gorm.IsRecordNotFoundError(err) {
		return application.NewErrNotFound(RECORD_NOT_FOUND)
	}

	return err
}

func (r *groupRepository) FindAlongRoles(g *models.Group) error {
	err := r.session.Preload(ROLES).First(g, g).Error
	if gorm.IsRecordNotFoundError(err) {
		return application.NewErrNotFound(RECORD_NOT_FOUND)
	}

	return err
}

func (r *groupRepository) FindAlongRolesAndUsers(g *models.Group) error {
	err := r.session.Preload(ROLES).Preload(USERS).First(g, g).Error
	if gorm.IsRecordNotFoundError(err) {
		return application.NewErrNotFound(RECORD_NOT_FOUND)
	}

	return err
}

func (r *groupRepository) Create(g *models.Group) error {
	err := r.session.Where(NAME, g.Name).First(&models.Group{}).Error
	if gorm.IsRecordNotFoundError(err) {
		return r.session.Set(GORM_AUTOUPDATE, false).Create(g).Error
	}

	if err == nil {
		return application.NewErrConflict(GROUP_ALREADY_EXIST)
	}

	return err
}

func (r *groupRepository) Delete(group *models.Group) error {
	if !r.session.Preload(USERS).First(&group, group).RecordNotFound() {
		if group.Name == ADMIN {
			return application.NewErrCustom(700, ADMIN_REMOVE_PERMISSION)
		}
		if len(group.Users) != 0 {
			return application.NewErrCustom(701, GROUP_HAS_USERS)
		}

		return r.session.Delete(&models.GroupsRole{}, GROUP_ID, group.ID).
			Delete(&models.Group{}, ID, group.ID).Error
	} else {
		return application.NewErrNotFound(RECORD_NOT_FOUND)
	}
}

func (r *groupRepository) Update(g *models.Group) error {
	alreadyExist := models.Group{}
	if !r.session.Where(NAME, g.Name).First(&alreadyExist).RecordNotFound() {
		if alreadyExist.Name == g.Name && alreadyExist.ID != g.ID {
			return application.NewErrConflict(GROUP_ALREADY_EXIST)
		}
	}

	group := models.Group{}
	if !r.session.Where(ID, g.ID).First(&group).RecordNotFound() {
		if group.Name == ADMIN {
			return application.NewErrCustom(700, ADMIN_UPDATE_PERMISSION)
		}
		if err := r.session.Unscoped().Delete(&models.GroupsRole{}, GROUP_ID, g.ID).Error; err != nil {
			return err
		}
		err := r.session.Set(GORM_AUTOUPDATE, false).Model(g).Updates(&g).Error
		return err
	} else {
		return application.NewErrNotFound(RECORD_NOT_FOUND)
	}
}

const (
	GORM_AUTOUPDATE            = "gorm:association_autoupdate"
	GORM_AUTOCREATE            = "gorm:association_autocreate"
	GROUP_ALREADY_EXIST        = "Group with this name already exist"
	ADMIN_REMOVE_PERMISSION    = "You can't remove Admin group"
	CUSTOMER_REMOVE_PERMISSION = "You can't remove Customer group"
	ADMIN_UPDATE_PERMISSION    = "You can't update Admin group"
	CUSTOMER_UPDATE_PERMISSION = "You can't update Customer group"
	GROUP_HAS_USERS            = "Group has users can't be removed"
	RECORD_NOT_FOUND           = "record not found"
	ADMIN                      = "مدیر"
	CUSTOMER                   = "Customer"
	ROLES                      = "Roles"
	USERS                      = "Users"
	ID                         = "id = ?"
	NAME                       = "name = ?"
	GROUP_ID                   = "group_id = ?"
)
