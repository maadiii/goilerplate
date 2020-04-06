package router

import (
	"goilerplate/usecase/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	ACCOUNTS        = "/accounts"
	ACCOUNTS_CREATE = ACCOUNTS + "/create"
	GROUPS          = "/groups"
	GROUPS_VIEW     = GROUPS + "/view/:id"
	GROUPS_CREATE   = GROUPS + "/create"
	GROUPS_UPDATE   = GROUPS + "/edit/:id"
	APIV1           = "/apiv1"
)

func Route(r *httprouter.Router, c controllers.IRootController) {
	handleView := c.GetBase().Handle
	handleRest := c.Apiv1().GetBase().Handle

	group := c.Groups().View
	groups := c.Groups().List
	addGroup := c.Groups().Add
	editGroup := c.Groups().Edit
	postGroup := c.Apiv1().Groups().Post
	putGroup := c.Apiv1().Groups().Put
	deleteGroup := c.Apiv1().Groups().Delete

	r.Handler(http.MethodGet, GROUPS, handleView(groups))
	r.Handler(http.MethodGet, GROUPS_VIEW, handleView(group))
	r.Handler(http.MethodGet, GROUPS_CREATE, handleView(addGroup))
	r.Handler(http.MethodGet, GROUPS_UPDATE, handleView(editGroup))
	r.Handler(http.MethodPost, APIV1+GROUPS, handleRest(postGroup))
	r.Handler(http.MethodDelete, APIV1+GROUPS, handleRest(deleteGroup))
	r.Handler(http.MethodPut, APIV1+GROUPS, handleRest(putGroup))

	addUser := c.Users().Add
	users := c.Users().List
	postUser := c.Apiv1().Users().Post

	r.Handler(http.MethodGet, ACCOUNTS_CREATE, handleView(addUser))
	r.Handler(http.MethodGet, ACCOUNTS, handleView(users))
	r.Handler(http.MethodPost, APIV1+ACCOUNTS, handleRest(postUser))
}
