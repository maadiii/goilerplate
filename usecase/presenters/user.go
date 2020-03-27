package presenters

import (
	"goilerplate/domain/models"

	"github.com/google/uuid"
)

type IUserPresenter interface {
	PresentSave(*models.User) UserPresent
}

type UserPresent struct {
	ID           uuid.UUID `json:"id"`
	MobileNumber string    `json:"mobileNumber"`
	FullName     string    `json:"fullName"`
}
