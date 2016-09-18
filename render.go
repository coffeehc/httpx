package web

import (
	"encoding/json"
	"io"

	"encoding/xml"
	"github.com/coffeehc/logger"
	"net/http"
)

var (
	DefaultCharset = "utf-8"

	Render_Json = render_Json{Charset: DefaultCharset}.Render
	Render_Xml  = render_Xml{Charset: DefaultCharset}.Render
	Render_Text = render_Text{Charset: DefaultCharset}.Render
)

func render(w io.Writer, data interface{}) error {
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

type Render func(w http.ResponseWriter, data interface{}) error

type render_Json struct {
	Charset string
}

func (this render_Json) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/json; charset="+this.Charset)
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return render(w, v)
}

type render_Text struct {
	Charset string
}

func (this render_Text) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "text/plain; charset="+this.Charset)
	return render(w, data)
}

type render_Xml struct {
	Charset string
}

func (this render_Xml) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/xml; charset="+this.Charset)
	v, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	return render(w, v)
}
