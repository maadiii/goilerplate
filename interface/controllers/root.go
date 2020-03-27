package controllers

import (
	"goilerplate/infrastructure/controller"
	"goilerplate/usecase/controllers"
)

type rootController struct {
	*controller.Controller
	Apiv1Controller controllers.IAPIV1Controller
	GroupController controllers.IGroupController
	UserController  controllers.IUserController
}

func NewRootController(c *controller.Controller) controllers.IRootController {
	return &rootController{Controller: c}
}

func (r *rootController) GetBase() *controller.Controller {
	return r.Controller
}

func (r *rootController) WithApiv1Controller(a controllers.IAPIV1Controller) controllers.IRootController {
	r.Apiv1Controller = a
	return r
}

func (r *rootController) Apiv1() controllers.IAPIV1Controller {
	return r.Apiv1Controller
}

func (r *rootController) WithGroupController(c controllers.IGroupController) controllers.IRootController {
	r.GroupController = c
	return r
}

func (r *rootController) Groups() controllers.IGroupController {
	return r.GroupController
}

func (r *rootController) WithUserController(c controllers.IUserController) controllers.IRootController {
	r.UserController = c
	return r
}

func (r *rootController) Users() controllers.IUserController {
	return r.UserController
}

type Apiv1Controller struct {
	*controller.RestController
	GroupController controllers.IGroupRestController
	UserController  controllers.IUserRestController
}

func NewApiv1Controller(c *controller.RestController) controllers.IAPIV1Controller {
	return &Apiv1Controller{RestController: c}
}

func (a Apiv1Controller) GetBase() *controller.RestController {
	return a.RestController
}

func (a Apiv1Controller) WithGroupController(c controllers.IGroupRestController) controllers.IAPIV1Controller {
	a.GroupController = c
	return a
}

func (a Apiv1Controller) Groups() controllers.IGroupRestController {
	return a.GroupController
}

func (a Apiv1Controller) WithUserController(c controllers.IUserRestController) controllers.IAPIV1Controller {
	a.UserController = c
	return a
}

func (a Apiv1Controller) Users() controllers.IUserRestController {
	return a.UserController
}
