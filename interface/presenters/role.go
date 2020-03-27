package presenters

import (
	"goilerplate/domain/models"
	up "goilerplate/usecase/presenters"
)

type rolePresenter struct{}

func NewRolePresenter() up.IRolePresenter {
	return &rolePresenter{}
}

func (presenter *rolePresenter) PresentAll(rls *[]models.Role) []up.RolePresent {
	roles := make([]up.RolePresent, len(*rls))
	for i, rl := range *rls {
		roles[i] = up.RolePresent{ID: rl.ID, FaName: rl.FaName, EnName: rl.EnName}
	}

	return roles
}
