package apiv1

import (
	"fmt"
	"goldfish/app"
	ctrl "goldfish/controllers"
	"goldfish/controllers/apiv1/admin"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ApiV1 struct {
	*ctrl.Controller
}

func New(c *ctrl.Controller) *ApiV1 {
	return &ApiV1{c}
}

func (a *ApiV1) Route(router *httprouter.Router) {
	g := admin.GroupController{a.Controller}
	u := admin.UserController{a.Controller}

	router.Handler(
		http.MethodGet,
		APPLICATION_URL,
		a.HandleRest(a.application),
	)

	router.Handler(http.MethodPost, ADMIN_ACCOUNTS, a.HandleRest(u.Post))

	router.Handler(http.MethodPost, ADMIN_GROUPS, a.HandleRest(g.Post))
	router.Handler(http.MethodDelete, ADMIN_GROUPS, a.HandleRest(g.Delete))
	router.Handler(http.MethodPatch, ADMIN_GROUPS, a.HandleRest(g.Patch))
}

func (a *ApiV1) application(ctx *app.Context) error {
	json := []byte(
		fmt.Sprintf(`{"title":"%v", "version":"%v"}`, app.Name, app.Version),
	)
	a.Response.Write(json)
	return nil
}
