package registry

import (
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/controller"
	"goilerplate/infrastructure/datastore"
	"goilerplate/interface/controllers"
	uc "goilerplate/usecase/controllers"
	"goilerplate/usecase/interactors"
	"goilerplate/usecase/presenters"
	"goilerplate/usecase/repositories"
)

type IRegistry interface {
	NewRootController() uc.IRootController

	// Group
	NewGroupController() uc.IGroupController
	NewGroupRestController() uc.IGroupRestController
	NewGroupInteractor() interactors.IGroupInteractor
	NewGroupRepository() repositories.IGroupRepository
	NewGroupPresenter() presenters.IGroupPresenter

	// Role
	NewRoleInteractor() interactors.IRoleInteractor
	NewRoleRepository() repositories.IRoleRepository
	NewRolePresenter() presenters.IRolePresenter

	// User
	NewUserController() uc.IUserController
	NewUserRestController() uc.IUserRestController
	NewUserInteractor() interactors.IUserInteractor
	NewUserRepository() repositories.IUserRepository
	NewUserPresenter() presenters.IUserPresenter
}

type registry struct {
	controller     *controller.Controller
	restController *controller.RestController
}

func NewRegistry() (IRegistry, error) {
	app, err := application.New()
	if err != nil {
		return nil, err
	}

	ctrl, err := controller.NewController(app)
	if err != nil {
		return nil, err
	}
	restController := controller.NewRestController(ctrl)

	return &registry{ctrl, restController}, nil
}

func NewTestRegistry() (IRegistry, error) {
	appConfig, err := application.InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := datastore.InitConfig()
	if err != nil {
		return nil, err
	}

	session, err := datastore.NewTestSession(dbConfig)
	if err != nil {
		return nil, err
	}

	app := &application.Application{appConfig, session}
	c, err := controller.NewController(app)
	if err != nil {
		return nil, err
	}
	rc := controller.NewRestController(c)

	return &registry{c, rc}, nil
}

func (r *registry) NewRootController() uc.IRootController {
	apiv1Controller := controllers.NewApiv1Controller(r.restController)
	apiv1Controller = apiv1Controller.WithGroupController(r.NewGroupRestController())
	apiv1Controller = apiv1Controller.WithUserController(r.NewUserRestController())

	ctrl := controllers.NewRootController(r.controller)
	ctrl = ctrl.WithApiv1Controller(apiv1Controller)
	ctrl = ctrl.WithGroupController(r.NewGroupController())
	ctrl = ctrl.WithUserController(r.NewUserController())

	return ctrl
}
