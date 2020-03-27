package controllers

import (
	"goilerplate/infrastructure/controller"
)

type IRootController interface {
	GetBase() *controller.Controller

	WithApiv1Controller(IAPIV1Controller) IRootController
	Apiv1() IAPIV1Controller

	WithGroupController(IGroupController) IRootController
	Groups() IGroupController

	WithUserController(IUserController) IRootController
	Users() IUserController
}

type IAPIV1Controller interface {
	GetBase() *controller.RestController

	WithGroupController(IGroupRestController) IAPIV1Controller
	Groups() IGroupRestController

	WithUserController(IUserRestController) IAPIV1Controller
	Users() IUserRestController
}
