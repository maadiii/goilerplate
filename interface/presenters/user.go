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

func (p *userPresenter) PresentAllAlongGroup(u *[]models.User) []up.UserAlongGroupPresent {
	users := make([]up.UserAlongGroupPresent, len(*u))
	for i, user := range *u {
		users[i] = up.UserAlongGroupPresent{
			ID:           user.ID,
			FullName:     user.FirstName + " " + user.LastName,
			MobileNumber: user.MobileNumber,
			GroupName:    user.Group.Name,
		}
	}

	return users
}

func (p *userPresenter) PresentCount(c *int) int {
	return *c
}

func PresentUser(u *models.User) up.UserPresent {
	return up.UserPresent{
		ID:           u.ID,
		FullName:     u.FirstName + " " + u.LastName,
		MobileNumber: u.MobileNumber,
	}
}
