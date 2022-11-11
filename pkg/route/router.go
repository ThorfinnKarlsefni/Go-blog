package route

import (
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
	route = r
}

func Name2URL(routeName string, pairs ...string) string {
	url, err := route.Get(routeName).URL(pairs...)

	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

func GetRouteVariable(parameteName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameteName]
}
