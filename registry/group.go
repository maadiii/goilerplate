package registry

import (
	"goilerplate/interface/controllers"
	"goilerplate/interface/presenters"
	"goilerplate/interface/repositories"
	uc "goilerplate/usecase/controllers"
	"goilerplate/usecase/interactors"
	up "goilerplate/usecase/presenters"
	ur "goilerplate/usecase/repositories"
)

func (r *registry) NewGroupController() uc.IGroupController {
	return controllers.NewGroupController(r.controller, r.NewGroupInteractor(), r.NewRoleInteractor())
}

func (r *registry) NewGroupRestController() uc.IGroupRestController {
	return controllers.NewGroupRestController(r.restController, r.NewGroupInteractor())
}

func (r *registry) NewGroupInteractor() interactors.IGroupInteractor {
	return interactors.NewGroupInteractor(r.NewGroupRepository(), r.NewGroupPresenter())
}

func (r *registry) NewGroupRepository() ur.IGroupRepository {
	return repositories.NewGroupRepository(r.controller.Application.DBSession)
}

func (r *registry) NewGroupPresenter() up.IGroupPresenter {
	return presenters.NewGroupPresenter()
}
