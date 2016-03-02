// pprof
package pprof

import (
	"github.com/coffeehc/web"
	"net/http/pprof"
)

func RegeditPprof(server *web.Server) {
	server.RegeditHttpHandlerFunc("/debug/pprof/{name}", web.GET, pprof.Index)
	server.RegeditHttpHandlerFunc("/debug/pprof/{name}", web.POST, pprof.Index)
	server.RegeditHttpHandlerFunc("/debug/pprof/cmdline", web.GET, pprof.Cmdline)
	server.RegeditHttpHandlerFunc("/debug/pprof/cmdline", web.POST, pprof.Cmdline)
	server.RegeditHttpHandlerFunc("/debug/pprof/profile", web.GET, pprof.Profile)
	server.RegeditHttpHandlerFunc("/debug/pprof/profile", web.POST, pprof.Profile)
	server.RegeditHttpHandlerFunc("/debug/pprof/symbol", web.GET, pprof.Symbol)
	server.RegeditHttpHandlerFunc("/debug/pprof/symbol", web.POST, pprof.Symbol)
	server.RegeditHttpHandlerFunc("/debug/pprof/trace", web.GET, pprof.Trace)
	server.RegeditHttpHandlerFunc("/debug/pprof/trace", web.POST, pprof.Trace)
}
