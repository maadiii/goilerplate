package interactors

import (
	"errors"
	"goilerplate/domain/models"
	"goilerplate/infrastructure/application"
	"goilerplate/infrastructure/utils"
	"goilerplate/usecase/presenters"
	"goilerplate/usecase/repositories"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type IUserInteractor interface {
	AllAlongGroup(int, string) ([]presenters.UserAlongGroupPresent, error)
	Save(*UserSave) (presenters.UserPresent, error)
	Count(string) (int, error)
}

type userInteractor struct {
	repository repositories.IUserRepository
	presenter  presenters.IUserPresenter
}

func NewUserInteractor(r repositories.IUserRepository, p presenters.IUserPresenter) IUserInteractor {
	return &userInteractor{r, p}
}

func (i *userInteractor) Save(u *UserSave) (presenters.UserPresent, error) {
	var userPresent presenters.UserPresent
	if err := u.Validate(); err != nil {
		return userPresent, application.NewErrValidation(err.Error())
	}

	user := models.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		MobileNumber: u.MobileNumber,
		Password:     []byte(u.Password),
		GroupID:      u.GroupID,
	}
	err := i.repository.Create(&user)
	if err != nil {
		return userPresent, err
	}

	userPresent = i.presenter.PresentSave(&user)
	return userPresent, nil
}

func (i *userInteractor) AllAlongGroup(p int, s string) ([]presenters.UserAlongGroupPresent, error) {
	users := []models.User{}
	err := i.repository.FindAllAlongGroup(&users, p, s)
	if err != nil {
		return nil, err
	}

	return i.presenter.PresentAllAlongGroup(&users), nil
}

func (i *userInteractor) Count(s string) (int, error) {
	var count int
	if err := i.repository.Count(s, &count); err != nil {
		return 0, err
	}

	return i.presenter.PresentCount(&count), nil
}

type UserSave struct {
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	MobileNumber string    `json:"mobileNumber"`
	Password     string    `json:"password"`
	GroupID      uuid.UUID `json:"groupId"`
}

func (u UserSave) Validate() error {
	err := validation.ValidateStruct(
		&u,
		validation.Field(
			&u.MobileNumber, validation.Required, validation.Length(11, 11), validation.Match(MOBILEREGEX),
		),
		validation.Field(&u.Password, validation.Required),
	)

	if err != nil {
		return err
	}

	if !utils.IsPasswordValid(u.Password) {
		return errors.New(INVALID_PASSWORD_FORMAT)
	}

	return nil
}

var (
	MOBILEREGEX = regexp.
		MustCompile("09(1[0-9]|3[1-9]|2[1-9])-?[0-9]{3}-?[0-9]{4}")
)

const (
	INVALID_PASSWORD_FORMAT = "password: must be in a valid format."
)
