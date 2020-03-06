package admin

import (
	"encoding/json"
	"goilerplate/app"
	"goilerplate/controllers"
	"goilerplate/domain/models"
	"goilerplate/domain/services"
	"goilerplate/types"

	"github.com/google/uuid"
)

type GroupController struct {
	*controllers.Controller
}

func (c GroupController) Post(ctx *app.Context) error {
	model := types.GroupAdd{}
	json.Unmarshal(ctx.Model, &model)
	if err := model.Validate(); err != nil {
		return app.NewErrValidation(err.Error())
	}

	group := &models.Group{
		Name:        model.Name,
		Description: model.Description,
	}

	for _, role := range model.Roles {
		group.Roles = append(group.Roles, models.Role{ID: role.ID})
	}

	return services.NewGroupService(ctx.DBSession).Save(group)
}

func (c GroupController) Patch(ctx *app.Context) error {
	model := types.GroupEdit{}
	json.Unmarshal(ctx.Model, &model)
	if err := model.Validate(); err != nil {
		return app.NewErrValidation(err.Error())
	}

	group := &models.Group{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
	}

	for _, role := range model.Roles {
		group.Roles = append(group.Roles, models.Role{ID: role.ID})
	}

	return services.NewGroupService(ctx.DBSession).Update(group)
}

func (g GroupController) Delete(ctx *app.Context) error {
	type Model struct {
		ID string `json:"id"`
	}

	model := Model{}
	json.Unmarshal(ctx.Model, &model)

	id, err := uuid.Parse(model.ID)
	if err != nil {
		return app.NewErrValidation(err.Error())
	}
	if id == uuid.Nil {
		return app.NewErrValidation("id: cannot be nil.")
	}

	return services.NewGroupService(ctx.DBSession).Delete(id)
}
