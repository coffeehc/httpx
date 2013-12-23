package web

import (
	"bufio"
	"errors"
	"html/template"
	"io"
)

type Tracsport interface {
	Out(writer io.Writer, reply *Reply) error
	ContentType() string
}

type text_Tracsport struct {
	Tracsport
}

func (this text_Tracsport) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (this text_Tracsport) Out(writer io.Writer, reply *Reply) error {
	switch reply.data.(type) {
	case string:
		_, err := writer.Write([]byte(reply.data.(string)))
		return err
	case []byte:
		_, err := writer.Write(reply.data.([]byte))
		return err
	case io.Reader:
		buf := bufio.NewReader(reply.data.(io.Reader))
		_, err := buf.WriteTo(writer)
		return err
	default:
		return errors.New("不接受的出书类型")
	}
}

type template_Tracsport struct {
	Tracsport
}

func (this template_Tracsport) ContentType() string {
	return "text/html; charset=utf-8"
}
func (this template_Tracsport) Out(writer io.Writer, reply *Reply) error {
	//TODO 模版缓存
	t, err := template.ParseFiles(reply.template)
	if err != nil {
		return err
	}
	return t.Execute(writer, reply.data)
}

type html_Tracsport struct {
	Tracsport
}

func (this html_Tracsport) ContentType() string {
	return "text/html"
}
func (this html_Tracsport) Out(writer io.Writer, reply *Reply) error {
	return Text.Out(writer, reply)
}

var Template template_Tracsport

var Text text_Tracsport

var Html html_Tracsport
