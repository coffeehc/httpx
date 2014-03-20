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
	return "application/json; charset=utf-8"
}

func (this text_Tracsport) ContentType() string {
	return "text/plain; charset=utf-8"
}

type html_Tracsport struct {
	Tracsport
}

func (this html_Tracsport) ContentType() string {
	return "text/html; charset=utf-8"
}

var Text text_Tracsport

var Html html_Tracsport

var Json json_Tracsport
