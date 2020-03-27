package controllers

import (
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/controller"
	"goilerplate/usecase/controllers"
	"goilerplate/usecase/interactors"
	views "goilerplate/views/admin"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shiyanhui/hero"
)

type groupController struct {
	*controller.Controller
	interactor     interactors.IGroupInteractor
	roleInteractor interactors.IRoleInteractor
}

func NewGroupController(
	c *controller.Controller, i interactors.IGroupInteractor, ri interactors.IRoleInteractor,
) controllers.IGroupController {
	return &groupController{c, i, ri}
}

func (c *groupController) Add(ctx *application.Context) error {
	roles, _ := c.roleInteractor.All()

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.AddGroup(roles, ctx.User, buffer)

	_, err := ctx.Response.Write(buffer.Bytes())

	return err
}

func (c *groupController) View(ctx *application.Context) error {
	params := httprouter.ParamsFromContext(ctx.Request.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return err
	}

	group, err := c.interactor.GetAlongRolesAndUsers(id)
	if err != nil {
		return err
	}

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.Group(group, ctx.User, buffer)
	_, err = ctx.Response.Write(buffer.Bytes())

	return err
}

func (c *groupController) Edit(ctx *application.Context) error {
	roles, _ := c.roleInteractor.All()

	params := httprouter.ParamsFromContext(ctx.Request.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return err
	}

	group, err := c.interactor.GetAlongSelectedRoles(id, &roles)
	if err != nil {
		return err
	}

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.EditGroup(group, ctx.User, buffer)
	_, err = ctx.Response.Write(buffer.Bytes())

	return err
}

func (c *groupController) List(ctx *application.Context) error {
	groups, _ := c.interactor.All()

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.GroupList(groups, ctx.User, buffer)
	_, err := ctx.Response.Write(buffer.Bytes())

	return err
}

type groupRestController struct {
	*controller.RestController
	interactor interactors.IGroupInteractor
}

func NewGroupRestController(c *controller.RestController, i interactors.IGroupInteractor) controllers.IGroupRestController {
	return &groupRestController{c, i}
}

func (c *groupRestController) Post(ctx *application.Context) error {
	var model interactors.GroupSave
	ctx.DecodeModel(&model)

	res, err := c.interactor.Save(&model)
	if err == nil {
		ctx.Json(res)
	}

	return err
}

func (c *groupRestController) Delete(ctx *application.Context) error {
	var model interactors.GroupDelete
	ctx.DecodeModel(&model)

	return c.interactor.Remove(&model)
}

func (c *groupRestController) Put(ctx *application.Context) error {
	var model interactors.GroupEdit
	ctx.DecodeModel(&model)

	res, err := c.interactor.Edit(&model)
	if err == nil {
		ctx.Json(res)
	}

	return err
}
