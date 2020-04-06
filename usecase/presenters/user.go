package presenters

import (
	"goilerplate/domain/models"

	"github.com/google/uuid"
)

type IUserPresenter interface {
	PresentSave(*models.User) UserPresent
	PresentAllAlongGroup(*[]models.User) []UserAlongGroupPresent
	PresentCount(*int) int
}

type UserPresent struct {
	ID           uuid.UUID `json:"id"`
	MobileNumber string    `json:"mobileNumber"`
	FullName     string    `json:"fullName"`
}

type UserAlongGroupPresent struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"fullName"`
	MobileNumber string    `json:"mobileNumber"`
	GroupName    string    `json:"groupName"`
}
