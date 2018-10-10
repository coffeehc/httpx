package httpx

import (
	"github.com/gin-gonic/gin/render"
)

type TemplateRender interface {
	JS(string, interface{}) render.Render
	CSS(string, interface{}) render.Render
	IMG(string, interface{}) render.Render
}
