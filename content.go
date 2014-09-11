// content
package web

import "inject"

const (
	CONTEXTTYPE_JSON       = "application/json"
	CONTEXTTYPE_JSON_RES   = "application/json; charset=utf-8"
	CONTEXTTYPE_PLAIN      = "text/plain"
	CONTEXTTYPE_PLAIN_RES  = "text/plain; charset=utf-8"
	CONTEXTTYPE_HTML       = "text/html"
	CONTEXTTYPE_HTML_RES   = "text/html; charset=utf-8"
	CONTEXTTYPE_XHTML      = "application/xhtml+xml"
	CONTEXTTYPE_XHTML_RES  = "application/xhtml+xml; charset=utf-8"
	CONTEXTTYPE_Binary     = "application/octet-stream"
	CONTEXTTYPE_Binary_RES = "application/octet-stream; charset=utf-8"
)

var (
	Bind_Key_Welcome        = "_welcomePath"
	Bind_Key_StaticResource = "_StaticRecouectFileSystem"
	Bind_Key_AppPath        = "_AppPath"
	Bind_Key_IsProjectModel = "_IsProjectModel"
	Bind_Key_TemplateDir    = "_TemplateDir"
	Bind_Key_Template       = "_Template"
	Bind_Key_Session        = inject.NameOf(inject.TypeOf((*Session)(nil)))
)

var default_charset = "utf-8"
