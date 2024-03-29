package middlewares

import (
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"net/http"
)

func Auth(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.Check() {
			flash.Warning("登陆用户才能访问此页面")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	}

}
