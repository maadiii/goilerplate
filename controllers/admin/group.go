package admin

import (
	"bytes"

	"goldfish/app"
	"goldfish/controllers"
	"goldfish/domain/models"
	"goldfish/domain/services"
	"goldfish/types"

	views "goldfish/views/admin"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shiyanhui/hero"
)

type GroupController struct {
	*controllers.Controller
}

func (c GroupController) Create(ctx *app.Context) error {
	roleList := []models.Role{}
	services.NewRoleService(ctx.DBSession).List(&roleList)

	buffer := new(bytes.Buffer)
	views.AddGroup(roleList, ctx.User, buffer)

	_, err := c.Response.Write(buffer.Bytes())

	return err
}

func (c GroupController) List(ctx *app.Context) error {
	groupList := []models.Group{}
	services.NewGroupService(ctx.DBSession).List(&groupList)

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.GroupList(groupList, ctx.User, buffer)

	_, err := c.Response.Write(buffer.Bytes())
	return err
}

func (c GroupController) Edit(ctx *app.Context) error {
	roles := []models.Role{}
	services.NewRoleService(ctx.DBSession).List(&roles)

	params := httprouter.ParamsFromContext(c.Request.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return err
	}

	group := models.Group{ID: id}
	err = services.NewGroupService(ctx.DBSession).SelectAndRoles(&group)
	if err != nil {
		return err
	}

	groupView := types.GroupWithSelectedRoles{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
	}
	for _, role := range roles {
		sRole := types.SelectedRole{
			ID:     role.ID,
			FaName: role.FaName,
			EnName: role.EnName,
		}
		for _, r := range group.Roles {
			if role.ID == r.ID {
				sRole.Selected = true
				break
			}
		}
		groupView.Roles = append(groupView.Roles, sRole)
	}

	buffer := hero.GetBuffer()
	defer hero.PutBuffer(buffer)
	views.UpdateGroup(groupView, ctx.User, buffer)

	_, err = c.Response.Write(buffer.Bytes())
	return err
}
