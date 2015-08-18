// pprof
package pprof

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/web"
)

func RegeditPprof(server *web.Server) {
	server.Regedit("/debug/pprof", web.GET, Index)
	server.Regedit("/debug/pprof/cmdline", web.GET, Cmdline)
	server.Regedit("/debug/pprof/profile", web.GET, Profile)
	server.Regedit("/debug/pprof/profile", web.POST, Profile)
	server.Regedit("/debug/pprof/symbol", web.GET, Symbol)
}

func Cmdline(request *http.Request, pathFragments map[string]string, reply *web.Reply) {
	reply.SetContentType("text/plain; charset=utf-8").With(strings.Join(os.Args, "\x00"))
}

func Profile(request *http.Request, pathFragments map[string]string, reply *web.Reply) {
	sec, _ := strconv.ParseInt(request.FormValue("seconds"), 10, 64)
	if sec == 0 {
		sec = 30
	}
	reply.SetContentType("application/octet-stream")
	r, w := io.Pipe()
	if err := pprof.StartCPUProfile(w); err != nil {
		reply.SetContentType("text/plain; charset=utf-8")
		reply.SetCode(http.StatusInternalServerError)
		reply.With(fmt.Sprintf("Could not enable CPU profiling: %s\n", err))
		return
	}
	go func() {
		time.Sleep(time.Duration(sec) * time.Second)
		pprof.StopCPUProfile()
		w.Close()
	}()
	reply.With(r)
}

// Symbol looks up the program counters listed in the request,
// responding with a table mapping program counters to function names.
// The package initialization registers it as /debug/pprof/symbol.
func Symbol(request *http.Request, pathFragments map[string]string, reply *web.Reply) {
	reply.SetContentType("text/plain; charset=utf-8")

	// We have to read the whole POST body before
	// writing any output.  Buffer the output here.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "num_symbols: 1\n")
	var b *bufio.Reader
	if request.Method == "POST" {
		b = bufio.NewReader(request.Body)
	} else {
		b = bufio.NewReader(strings.NewReader(request.URL.RawQuery))
	}
	for {
		word, err := b.ReadSlice('+')
		if err == nil {
			word = word[0 : len(word)-1] // trim +
		}
		pc, _ := strconv.ParseUint(string(word), 0, 64)
		if pc != 0 {
			f := runtime.FuncForPC(uintptr(pc))
			if f != nil {
				fmt.Fprintf(&buf, "%#x %s\n", pc, f.Name())
			}
		}

		// Wait until here to check for err; the last
		// symbol will have an err because it doesn't end in +.
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(&buf, "reading request: %v\n", err)
			}
			break
		}
	}
	reply.With(string(buf.Bytes()))
}

type handler string

func (name handler) RequestHandler(request *http.Request, pathFragments map[string]string, reply *web.Reply) {
	reply.SetContentType("text/plain; charset=utf-8")
	debug, _ := strconv.Atoi(request.FormValue("debug"))
	p := pprof.Lookup(string(name))
	if p == nil {
		reply.SetCode(404).With(fmt.Sprintf("Unknown profile: %s\n", name))
		return
	}
	gc, _ := strconv.Atoi(request.FormValue("gc"))
	if name == "heap" && gc > 0 {
		runtime.GC()
	}
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		p.WriteTo(w, debug)
	}()
	reply.With(r)
}

// Index responds with the pprof-formatted profile named by the request.
// For example, "/debug/pprof/heap" serves the "heap" profile.
// Index responds to a request for "/debug/pprof/" with an HTML page
// listing the available profiles.
func Index(request *http.Request, pathFragments map[string]string, reply *web.Reply) {
	if strings.HasPrefix(request.URL.Path, "/debug/pprof/") {
		name := strings.TrimPrefix(request.URL.Path, "/debug/pprof/")
		if name != "" {
			handler(name).RequestHandler(request, pathFragments, reply)
			return
		}
	}
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		profiles := pprof.Profiles()
		if err := indexTmpl.Execute(w, profiles); err != nil {
			logger.Error("出现错误:%s", err)
		}
	}()
	reply.With(r)
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
<head>
<title>/debug/pprof/</title>
</head>
/debug/pprof/<br>
<br>
<body>
profiles:<br>
<table>
{{range .}}
<tr><td align=right>{{.Count}}<td><a href="/debug/pprof/{{.Name}}?debug=1">{{.Name}}</a>
{{end}}
</table>
<br>
<a href="/debug/pprof/goroutine?debug=2">full goroutine stack dump</a><br>
</body>
</html>
`))
