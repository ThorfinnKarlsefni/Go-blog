package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

func Render(w io.Writer, data interface{}, tplFiels ...string) {
	viewDir := "resources/views/"

	for i, f := range tplFiels {
		tplFiels[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")

	allFiles := append(layoutFiles, tplFiels...)

	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)

	logger.LogError(err)

	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
