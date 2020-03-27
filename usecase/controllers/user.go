package controllers

import "goilerplate/infrastructure/application"

type IUserController interface {
	Add(*application.Context) error
}

type IUserRestController interface {
	Post(*application.Context) error
}
