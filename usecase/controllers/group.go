package controllers

import "goilerplate/infrastructure/application"

type IGroupController interface {
	Add(*application.Context) error
	Edit(*application.Context) error
	List(*application.Context) error
  View(*application.Context) error
}

type IGroupRestController interface {
	Post(*application.Context) error
	Delete(*application.Context) error
	Put(*application.Context) error
}
