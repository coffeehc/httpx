package pprof

import (
	"net/http/pprof"

	"git.xiagaogao.com/coffee/httpx"
)

//RegeditPprof register pprof to http server
func RegeditPprof(server httpx.Server) {
	server.RegisterHandlerFunc("/debug/pprof/{name}", httpx.GET, pprof.Index)
	server.RegisterHandlerFunc("/debug/pprof/{name}", httpx.POST, pprof.Index)
	server.RegisterHandlerFunc("/debug/pprof/cmdline", httpx.GET, pprof.Cmdline)
	server.RegisterHandlerFunc("/debug/pprof/cmdline", httpx.POST, pprof.Cmdline)
	server.RegisterHandlerFunc("/debug/pprof/profile", httpx.GET, pprof.Profile)
	server.RegisterHandlerFunc("/debug/pprof/profile", httpx.POST, pprof.Profile)
	server.RegisterHandlerFunc("/debug/pprof/symbol", httpx.GET, pprof.Symbol)
	server.RegisterHandlerFunc("/debug/pprof/symbol", httpx.POST, pprof.Symbol)
	server.RegisterHandlerFunc("/debug/pprof/trace", httpx.GET, pprof.Trace)
	server.RegisterHandlerFunc("/debug/pprof/trace", httpx.POST, pprof.Trace)
}
