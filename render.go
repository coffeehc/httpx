package web

import (
	"io"

	"encoding/xml"
	"github.com/coffeehc/logger"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
)

var (
	DefaultCharset = "utf-8"

	Default_Render_Json = Render_Json{Charset: DefaultCharset}
	Default_Render_Xml  = Render_Xml{Charset: DefaultCharset}
	Default_Render_Text = Render_Text{Charset: DefaultCharset}
)

type Render interface {
	ContentType() string
	Write(w http.ResponseWriter, data interface{}) error
}

func renderReply(w http.ResponseWriter, render Render, data interface{}) error {
	contentType := render.ContentType()
	if contentType == "" {
		contentType = "text/plain"
	}
	w.Header().Set("content-type", contentType)
	return render.Write(w, data)

}

func writeBody(w io.Writer, data interface{}) error {
	var err error = nil
	switch v := data.(type) {
	case string:
		_, err = w.Write([]byte(v))
	case *string:
		_, err = w.Write([]byte(*v))
	case []byte:
		_, err = w.Write(v)
	case io.Reader:
		_, err = io.Copy(w, v)
		defer func() {
			if closer, ok := v.(io.Closer); ok {
				closer.Close()
			}
		}()
	case nil:
		logger.Warn("不返回任何信息")
	}
	return err
}

type Render_Json struct {
	Charset string
}

func (this Render_Json) ContentType() string {
	return "application/json; charset=" + this.Charset
}

func (this Render_Json) Write(w http.ResponseWriter, data interface{}) error {
	v, err := ffjson.Marshal(data)
	if err != nil {
		return err
	}
	return writeBody(w, v)
}

type Render_Text struct {
	Charset string
}

func (this Render_Text) ContentType() string {
	return "text/plain; charset=" + this.Charset
}

func (this Render_Text) Write(w http.ResponseWriter, data interface{}) error {
	return writeBody(w, data)
}

type Render_Xml struct {
	Charset string
}

func (this Render_Xml) ContentType() string {
	return "text/plain; charset=" + this.Charset
}

func (this Render_Xml) Write(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/xml; charset="+this.Charset)
	v, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	return writeBody(w, v)
}
