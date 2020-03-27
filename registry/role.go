package registry

import (
	"goilerplate/interface/presenters"
	"goilerplate/interface/repositories"
	"goilerplate/usecase/interactors"
	up "goilerplate/usecase/presenters"
	ur "goilerplate/usecase/repositories"
)

func (r *registry) NewRoleInteractor() interactors.IRoleInteractor {
	return interactors.NewRoleInteractor(r.NewRoleRepository(), r.NewRolePresenter())
}

func (r *registry) NewRoleRepository() ur.IRoleRepository {
	return repositories.NewRoleRepository(r.controller.Application.DBSession)
}

func (r *registry) NewRolePresenter() up.IRolePresenter {
	return presenters.NewRolePresenter()
}
