package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	Password     []byte     `gorm:"not null"`
	MobileNumber string     `gorm:"type:varchar(11);unique_index""`
	FirstName    string     `gorm:"type:varchar(64);index"`
	LastName     string     `gorm:"type:varchar(64);index"`
	IsActive     bool       `gorm:"not null"`
	GroupID      uuid.UUID  `gorm:"type:uuid;not null"`
	Group        Group
}

type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;priamry_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Name        string     `gorm:"type:varchar(64);not null;unique;"`
	Description string     `gorm:"type:varchar(256);not null"`

	Roles []Role `gorm:"many2many:groups_roles"`
	Users []User
}

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;priamry_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	FaName    string     `gorm:"not null; type:varchar(64);unique"`
	EnName    string     `gorm:"not null; type:varchar(64);unique"`

	Groups []Group `gorm:"many2many:groups_roles"`
}

type GroupsRole struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	GroupID   uuid.UUID  `gorm:"type:uuid;primary_key"`
	RoleID    uuid.UUID  `grom:"type:uuid;primary_key"`

	Group Group
	Role  Role
}

func (u *User) BeforeCreate(scope *gorm.Scope) (err error) {
	PRE_ERROR := "domain/model/user.go BeforeCreate(),"
	hashed, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf(
			"%s %s %v",
			PRE_ERROR,
			"Failed to password hashing:",
			err,
		)
	}

	err = scope.SetColumn("Password", hashed)
	if err != nil {
		return err
	}

	return scope.SetColumn("ID", uuid.New())
}

func (g *Group) BeforeCreate(scope *gorm.Scope) (err error) {
	return scope.SetColumn("ID", uuid.New())
}

func (r *Role) BeforeCreate(scope *gorm.Scope) (err error) {
	return scope.SetColumn("ID", uuid.New())
}
