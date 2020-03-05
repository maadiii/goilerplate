package types

import (
	"errors"

	"github.com/google/uuid"
)

type Role struct {
	EnName string    `json:"enName"`
	FaName string    `json:"faName"`
	ID     uuid.UUID `json:"id"`
}

type SelectedRole struct {
	ID       uuid.UUID
	FaName   string
	EnName   string
	Selected bool
}

func (r Role) Validate() error {
	if r.ID != uuid.Nil {
		return nil
	}
	return errors.New("invalid role")
}
