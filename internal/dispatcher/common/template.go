package common

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

var FuncMap = template.FuncMap{
	"Marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
}

// Template
type Template struct {
	tpl *template.Template
}

// Render
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tpl.ExecuteTemplate(w, name, data)
}

// NewTemplate
func NewTemplate(tpl *template.Template) *Template {
	return &Template{tpl: tpl}
}
