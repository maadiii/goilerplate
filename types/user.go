package types

import (
	"errors"
	"goilerplate/app"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

var (
	MOBILEREGEX = regexp.
		MustCompile("09(1[0-9]|3[1-9]|2[1-9])-?[0-9]{3}-?[0-9]{4}")
)

type User struct {
	ID           uint64 `json:"id"`
	MobileNumber string `json:"mobileNumber"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type UserAdd struct {
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	MobileNumber string    `json:"mobileNumber"`
	Password     string    `json:"password"`
	GroupID      uuid.UUID `json:"groupId"`
}

type UserForJWT struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Roles     []string
}

func (u UserAdd) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(
			&u.MobileNumber,
			validation.Required,
			validation.Length(11, 11),
			validation.Match(MOBILEREGEX),
		),
		validation.Field(
			&u.Password,
			validation.Required,
		),
	)

	if err != nil {
		return err
	}

	if !app.IsPasswordValid(u.Password) {
		return errors.New(INVALID_PASSWORD_FORMAT)
	}

	return nil
}
