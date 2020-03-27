package controllers

import (
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/controller"
	"goilerplate/usecase/controllers"
	"goilerplate/usecase/interactors"
	views "goilerplate/views/admin"

	"github.com/shiyanhui/hero"
)

type userController struct {
	*controller.Controller
	Interactor      interactors.IUserInteractor
	GroupInteractor interactors.IGroupInteractor
}

func NewUserController(
	c *controller.Controller, i interactors.IUserInteractor, g interactors.IGroupInteractor,
) controllers.IUserController {
	return &userController{c, i, g}
}

func (userController *userController) Add(context *application.Context) error {
	groups, _ := userController.GroupInteractor.All()

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)

	views.AddUser(groups, context.User, buffer)
	_, err := context.Response.Write(buffer.Bytes())

	return err
}

type userRestController struct {
	*controller.RestController
	interactor interactors.IUserInteractor
}

func NewUserRestController(
	c *controller.RestController, i interactors.IUserInteractor,
) controllers.IUserRestController {
	return &userRestController{c, i}
}

func (c *userRestController) Post(ctx *application.Context) error {
	var model interactors.UserSave
	ctx.DecodeModel(&model)

	res, err := c.interactor.Save(&model)
	if err == nil {
		ctx.Json(res)
	}

	return err
}
