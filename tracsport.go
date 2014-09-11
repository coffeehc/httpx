package web

type Tracsport interface {
	ContentType() string
}

type text_Tracsport struct {
	Tracsport
}

type json_Tracsport struct {
	Tracsport
}

func (this json_Tracsport) ContentType() string {
	return CONTEXTTYPE_JSON_RES
}

func (this text_Tracsport) ContentType() string {
	return CONTEXTTYPE_PLAIN_RES
}

type html_Tracsport struct {
	Tracsport
}

func (this html_Tracsport) ContentType() string {
	return CONTEXTTYPE_HTML_RES
}

var Text text_Tracsport

var Html html_Tracsport

var Json json_Tracsport
