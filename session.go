// session
package web

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/coffeehc/logger"
)

const (
	SESSION_COOKIE_ID = "_sid_"
)

type SessionStoreConf struct {
	MaxValidTime time.Duration
	Gclifetime   time.Duration
}

type SessionStore interface {
	GetSession(req *http.Request, w http.ResponseWriter) Session
	SaveSession(session Session)
	RemoveSession(session Session)
	GC()
	GetGclifetime() time.Duration
	InitSessionStore(conf SessionStoreConf)
}

type Session interface {
	GetCreationTime() time.Time
	GetId() string
	Get(key string) interface{}        //获取Session值
	Set(key string, value interface{}) //设置Session值
	Remove(key string)
	GetLastAccessedTime() time.Time
	UpdateLastAccessedTime()
	Reset()
}

type MemSessionStore struct {
	lock         *sync.RWMutex
	maxValidTime time.Duration
	sessions     map[string]*list.Element
	list         *list.List
	gcLifeTime   time.Duration
}

func (this *MemSessionStore) InitSessionStore(conf SessionStoreConf) {
	this.sessions = make(map[string]*list.Element)
	this.lock = new(sync.RWMutex)
	this.list = new(list.List)
	this.maxValidTime = conf.MaxValidTime
	this.gcLifeTime = conf.Gclifetime

}

func (this *MemSessionStore) GetGclifetime() time.Duration {
	return this.gcLifeTime
}
func (this *MemSessionStore) GetSession(req *http.Request, w http.ResponseWriter) Session {
	var sessionId string
	cookie, err := req.Cookie(SESSION_COOKIE_ID)
	if err != nil || cookie.Value == "" {
		sessionId = NewSessionId()
		cookie = &http.Cookie{Name: SESSION_COOKIE_ID,
			Value:    sessionId,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			MaxAge:   int(this.maxValidTime / time.Second),
		}
	} else {
		sessionId = cookie.Value
	}
	http.SetCookie(w, cookie)
	req.AddCookie(cookie)
	this.lock.RLock()
	if element, ok := this.sessions[sessionId]; ok {
		element.Value.(Session).UpdateLastAccessedTime()
		this.lock.RUnlock()
		return element.Value.(Session)
	} else {
		this.lock.RUnlock()
		this.lock.Lock()
		now := time.Now()
		newsess := &SessionImpl{Id: sessionId, CreatTime: now, LastAccessedTime: now, Data: make(map[string]interface{}), lock: new(sync.RWMutex)}
		element := this.list.PushBack(newsess)
		this.sessions[sessionId] = element
		this.lock.Unlock()
		logger.Debug("创建一个新的Session:%s", sessionId)
		return newsess
	}
	return nil
}

func (this *MemSessionStore) RemoveSession(session Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if element, ok := this.sessions[session.GetId()]; ok {
		delete(this.sessions, session.GetId())
		this.list.Remove(element)
	}
}

func (this *MemSessionStore) SaveSession(session Session) {
	//内存操作，不需要存储
}

func (this *MemSessionStore) GC() {
	this.lock.RLock()
	for {
		element := this.list.Back()
		if element == nil {
			break
		}
		session := element.Value.(Session)
		if session.GetLastAccessedTime().Sub(session.GetCreationTime()) > this.maxValidTime {
			this.lock.RUnlock()
			this.lock.Lock()
			this.list.Remove(element)
			delete(this.sessions, session.GetId())
			this.lock.Unlock()
			this.lock.RLock()
		} else {
			break
		}
	}
	this.lock.RUnlock()
}

type SessionImpl struct {
	lock             *sync.RWMutex
	Id               string
	Data             map[string]interface{}
	CreatTime        time.Time
	LastAccessedTime time.Time
}

func (this *SessionImpl) Reset() {
	this.Data = make(map[string]interface{})
}
func (this *SessionImpl) GetCreationTime() time.Time {
	return this.CreatTime
}
func (this *SessionImpl) GetId() string {
	return this.Id
}
func (this *SessionImpl) Get(key string) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.Data[key]
}
func (this *SessionImpl) Set(key string, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Data[key] = value
}
func (this *SessionImpl) Remove(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.Data, key)
}
func (this *SessionImpl) GetLastAccessedTime() time.Time {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.LastAccessedTime
}
func (this *SessionImpl) UpdateLastAccessedTime() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.LastAccessedTime = time.Now()
}

func NewSessionId() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return NewSessionId()
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
