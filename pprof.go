package httpx

//import (
//	"github.com/gofiber/fiber/v2"
//	"net/http/pprof"
//)
//
//func RegisterPprof(route fiber.Router) {
//	route = route.Group("/debug/pprof")
//	route.Get("/", indexHandler())
//	route.Get("/heap", heapHandler())
//	route.Get("/goroutine", goroutineHandler())
//	route.Get("/block", blockHandler())
//	route.Get("/threadcreate", threadCreateHandler())
//	route.Get("/cmdline", cmdlineHandler())
//	route.Get("/profile", profileHandler())
//	route.Get("/symbol", symbolHandler())
//	route.Get("/symbol", symbolHandler())
//	route.Get("/trace", traceHandler())
//	route.Get("/mutex", mutexHandler())
//	route.Get("/allocs", allocsHandler())
//
//}
//
//// IndexHandler will pass the call from /debug/pprof to pprof
//func indexHandler() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		pprof.Index(ctx.Writer, ctx.Request)
//	}
//}
//
//// HeapHandler will pass the call from /debug/pprof/heap to pprof
//func heapHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("heap"))
//}
//
//// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
//func allocsHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("allocs"))
//}
//
//// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
//func goroutineHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("goroutine"))
//}
//
//// BlockHandler will pass the call from /debug/pprof/block to pprof
//func blockHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("block"))
//}
//
//// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
//func threadCreateHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("threadcreate"))
//}
//
//// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
//func cmdlineHandler() fiber.Handler {
//	return func(ctx *gin.Context) {
//		pprof.Cmdline(ctx.Writer, ctx.Request)
//	}
//}
//
//// ProfileHandler will pass the call from /debug/pprof/profile to pprof
//func profileHandler() fiber.Handler {
//	return func(ctx *gin.Context) {
//		pprof.Profile(ctx.Writer, ctx.Request)
//	}
//}
//
//// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
//func symbolHandler() fiber.Handler {
//	return func(ctx *gin.Context) {
//		pprof.Symbol(ctx.Writer, ctx.Request)
//	}
//}
//
//// TraceHandler will pass the call from /debug/pprof/trace to pprof
//func traceHandler() fiber.Handler {
//	return func(ctx *gin.Context) {
//		pprof.Trace(ctx.Writer, ctx.Request)
//	}
//}
//
//// MutexHandler will pass the call from /debug/pprof/mutex to pprof
//func mutexHandler() fiber.Handler {
//	return gin.WrapH(pprof.Handler("mutex"))
//}
