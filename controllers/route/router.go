package route

import (
	. "goilerplate/controllers"
	"goilerplate/controllers/admin"
	"goilerplate/controllers/apiv1"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Route(c *Controller, r *httprouter.Router) {

	// Apiv1 routes
	apiv1 := apiv1.New(c)
	apiv1.Route(r)

	g := admin.GroupController{c}
	u := admin.UserController{c}

	r.Handler(http.MethodGet, ADMIN_GROUPS_CREATE, c.HandleView(g.Create))
	r.Handler(http.MethodGet, ADMIN_GROUPS_EDIT, c.HandleView(g.Edit))
	r.Handler(http.MethodGet, ADMIN_GROUPS, c.HandleView(g.List))

	r.Handler(http.MethodGet, ADMIN_ACCOUNTS_CREATE, c.HandleView(u.Create))
}
