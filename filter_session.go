// session
package web

import (
	"inject"
	"net/http"
	"time"
)

type SessionFilter struct {
	Pattern      string
	sessionStore SessionStore
	reponseName  string
}

func (this *SessionFilter) GetPattern() string {
	return this.Pattern
}

func (this *SessionFilter) Init(conf *WebConfig) {
	if store, ok := conf.GetInterface(inject.NameOf(inject.TypeOf((*SessionStore)(nil)))).(SessionStore); ok {
		this.sessionStore = store
	} else {
		this.sessionStore = new(MemSessionStore)
		this.sessionStore.InitSessionStore(SessionStoreConf{MaxValidTime: 30 * time.Minute, Gclifetime: time.Minute})
	}
	this.reponseName = inject.NameOf(inject.TypeOf((*http.ResponseWriter)(nil)))
	go func() {
		this.sessionStore.GC()
		time.AfterFunc(this.sessionStore.GetGclifetime(), func() { this.sessionStore.GC() })
	}()
}

func (this *SessionFilter) DoFilter(req *http.Request, reply *Reply, chain FilterChain) {
	writer := reply.GetInterface(this.reponseName).(http.ResponseWriter)
	session := this.sessionStore.GetSession(req, writer)
	reply.Binding(session, nil, Bind_Key_Session)
	chain.DoFilter(req, reply)
	this.sessionStore.SaveSession(session)

}

func (this *SessionFilter) Destroy() {
}
