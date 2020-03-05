package admin

import (
	"bytes"
	"goldfish/app"
	"goldfish/controllers"
	"goldfish/domain/models"
	"goldfish/domain/services"
	views "goldfish/views/admin"
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
