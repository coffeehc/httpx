package render

import (
	"html/template"
)

func LoadTemplateGlob(pattern string) *template.Template{
	return template.Must(template.New("").Delims("{{", "}}").Funcs(template.FuncMap{}).ParseGlob(pattern))
}
