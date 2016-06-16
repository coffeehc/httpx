// pprof
package pprof

import (
	"net/http/pprof"

	"github.com/coffeehc/web"
)

func RegeditPprof(server *web.Server) {
	server.RegisterHttpHandlerFunc("/debug/pprof/{name}", web.GET, pprof.Index)
	server.RegisterHttpHandlerFunc("/debug/pprof/{name}", web.POST, pprof.Index)
	server.RegisterHttpHandlerFunc("/debug/pprof/cmdline", web.GET, pprof.Cmdline)
	server.RegisterHttpHandlerFunc("/debug/pprof/cmdline", web.POST, pprof.Cmdline)
	server.RegisterHttpHandlerFunc("/debug/pprof/profile", web.GET, pprof.Profile)
	server.RegisterHttpHandlerFunc("/debug/pprof/profile", web.POST, pprof.Profile)
	server.RegisterHttpHandlerFunc("/debug/pprof/symbol", web.GET, pprof.Symbol)
	server.RegisterHttpHandlerFunc("/debug/pprof/symbol", web.POST, pprof.Symbol)
	server.RegisterHttpHandlerFunc("/debug/pprof/trace", web.GET, pprof.Trace)
	server.RegisterHttpHandlerFunc("/debug/pprof/trace", web.POST, pprof.Trace)
}
