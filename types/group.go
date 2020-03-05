package types

import (
	"goldfish/app"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Roles       []Role    `json:"roles"`
}

type GroupAdd struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Roles       []Role `json:"roles"`
}

type GroupEdit struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Roles       []Role    `json:"roles"`
}

func (g GroupAdd) Validate() error {
	return validation.ValidateStruct(
		&g,
		validation.Field(
			&g.Name,
			validation.Required,
			validation.Length(6, 128),
		),
		validation.Field(
			&g.Description,
			validation.Required,
			validation.Length(32, 256),
		),
		validation.Field(
			&g.Roles,
			validation.Required,
			validation.Length(1, 100),
		),
	)
}

func (g GroupEdit) Validate() error {
	return validation.ValidateStruct(
		&g,
		validation.Field(
			&g.ID,
			validation.Required,
			validation.Length(36, 36),
			validation.By(app.CheckUUID),
		),
		validation.Field(
			&g.Name,
			validation.Required,
			validation.Length(6, 128),
		),
		validation.Field(
			&g.Description,
			validation.Required,
			validation.Length(32, 256),
		),
		validation.Field(&g.Roles, validation.Required),
	)
}

type GroupWithSelectedRoles struct {
	ID          uuid.UUID
	Name        string
	Description string
	Roles       []SelectedRole
}
