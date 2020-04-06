package controllers

import (
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/controller"
	"goilerplate/usecase/controllers"
	"goilerplate/usecase/interactors"
	views "goilerplate/views/admin"
	"strconv"

	"github.com/shiyanhui/hero"
)

type userController struct {
	*controller.Controller
	interactor      interactors.IUserInteractor
	groupInteractor interactors.IGroupInteractor
}

func NewUserController(
	c *controller.Controller, i interactors.IUserInteractor, g interactors.IGroupInteractor,
) controllers.IUserController {
	return &userController{c, i, g}
}

func (c *userController) Add(ctx *application.Context) error {
	groups, _ := c.groupInteractor.All()

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)

	views.AddUser(groups, ctx.User, buffer)
	_, err := ctx.Response.Write(buffer.Bytes())

	return err
}

func (c *userController) List(ctx *application.Context) error {
	pageNumber, _ := strconv.Atoi(ctx.Request.URL.Query().Get(PAGE_NUMBER))
	search := ctx.Request.URL.Query().Get(SEARCH)
	users, err := c.interactor.AllAlongGroup(pageNumber, search)
	if err != nil {
		return err
	}

	usersCount, _ := c.interactor.Count(search)
	lastPage := usersCount / 10
	if usersCount%10 != 0 {
		lastPage++
	}

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)

	views.UsersList(users, pageNumber, lastPage, search, ctx.User, buffer)
	_, err = ctx.Response.Write(buffer.Bytes())

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

const (
	PAGE_NUMBER = "page"
	SEARCH      = "search"
)
