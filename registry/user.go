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

func (r *registry) NewUserController() uc.IUserController {
	return controllers.NewUserController(r.controller, r.NewUserInteractor(), r.NewGroupInteractor())
}

func (r *registry) NewUserRestController() uc.IUserRestController {
	return controllers.NewUserRestController(r.restController, r.NewUserInteractor())
}

func (r *registry) NewUserInteractor() interactors.IUserInteractor {
	return interactors.NewUserInteractor(r.NewUserRepository(), r.NewUserPresenter())
}

func (r *registry) NewUserRepository() ur.IUserRepository {
	return repositories.NewUserRepository(r.controller.Application.DBSession)
}

func (r *registry) NewUserPresenter() up.IUserPresenter {
	return presenters.NewUserPresenter()
}
