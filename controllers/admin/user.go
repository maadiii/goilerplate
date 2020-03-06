package admin

import (
	"bytes"
	"goilerplate/app"
	"goilerplate/controllers"
	"goilerplate/domain/models"
	"goilerplate/domain/services"
	views "goilerplate/views/admin"
)

type UserController struct {
	*controllers.Controller
}

func (c UserController) Create(ctx *app.Context) error {
	groupList := []models.Group{}
	services.NewGroupService(ctx.DBSession).List(&groupList)

	buffer := new(bytes.Buffer)
	views.AddUser(groupList, ctx.User, buffer)

	_, err := c.Response.Write(buffer.Bytes())
	return err
}
