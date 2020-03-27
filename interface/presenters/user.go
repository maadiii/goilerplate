package presenters

import (
	"goilerplate/domain/models"
	up "goilerplate/usecase/presenters"
)

type userPresenter struct{}

func NewUserPresenter() up.IUserPresenter {
	return &userPresenter{}
}

func (p *userPresenter) PresentSave(u *models.User) up.UserPresent {
	return PresentUser(u)
}

func PresentUser(u *models.User) up.UserPresent {
	return up.UserPresent{
		ID:           u.ID,
		FullName:     u.FirstName + " " + u.LastName,
		MobileNumber: u.MobileNumber,
	}
}
