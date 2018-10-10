package render

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/render"
)

type TemplateRender interface {
	Data(contextType string, name string, data interface{}) render.Render
	HTML(name string, data interface{}) render.Render
	JS(name string, data interface{}) render.Render
	CSS(name string, data interface{}) render.Render
}

func NewTemplateRender(tmpl *template.Template) TemplateRender {
	return &templateRenderImpl{
		Template: tmpl,
	}
}

type templateRenderImpl struct {
	Template *template.Template
}

func (t *templateRenderImpl) Data(contextType string, name string, data interface{}) render.Render {
	return &templateRender{
		Template:    t.Template,
		Name:        name,
		Data:        data,
		ContextType: []string{contextType},
	}
}

func (t *templateRenderImpl) HTML(name string, data interface{}) render.Render {
	return t.Data(htmlContentType, name, data)
}

func (t *templateRenderImpl) JS(name string, data interface{}) render.Render {
	return t.Data(jsContentType, name, data)
}

func (t *templateRenderImpl) CSS(name string, data interface{}) render.Render {
	return t.Data(cssContentType, name, data)
}

type templateRender struct {
	Template    *template.Template
	Name        string
	Data        interface{}
	ContextType []string
}

func (r templateRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

func (r templateRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, r.ContextType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
