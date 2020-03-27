package interactors

import (
	"goilerplate/domain/models"
	"goilerplate/usecase/presenters"
	"goilerplate/usecase/repositories"
)

type IRoleInteractor interface {
	All() ([]presenters.RolePresent, error)
}

type roleInteractor struct {
	repository repositories.IRoleRepository
	presenter  presenters.IRolePresenter
}

func NewRoleInteractor(r repositories.IRoleRepository, p presenters.IRolePresenter) IRoleInteractor {
	return &roleInteractor{r, p}
}

func (interactor *roleInteractor) All() ([]presenters.RolePresent, error) {
	roles := []models.Role{}
	err := interactor.repository.FindAll(&roles)
	if err != nil {
		return nil, err
	}

	return interactor.presenter.PresentAll(&roles), nil
}
