package knot

import (
	"sync"
	"time"

	"github.com/eaciit/toolkit"
)

type Sessions map[string]toolkit.M

var (
	sessionCookieId string
	sessions        Sessions      = make(map[string]toolkit.M)
	sessionLocker   *sync.RWMutex = new(sync.RWMutex)
)

func SetSessionCookieId(id string) {
	sessionCookieId = id
}

func SessionCookieId() string {
	if sessionCookieId == "" {
		sessionCookieId = "KnotSessionId"
	}
	return sessionCookieId
}

func (s Sessions) InitTokenBucket(tokenid string) {
	sessionLocker.Lock()
	if _, b := s[tokenid]; !b {
		s[tokenid] = toolkit.M{}
	}
	sessionLocker.Unlock()
}

func (s Sessions) Set(tokenid, key string, value interface{}) {
	s.InitTokenBucket(tokenid)

	sessionLocker.Lock()
	s[tokenid].Set(key, value)
	sessionLocker.Unlock()
}

func (s Sessions) Get(tokenid, key string, def interface{}) interface{} {
	s.InitTokenBucket(tokenid)

	sessionLocker.RLock()
	value := s[tokenid].Get(key, def)
	sessionLocker.RUnlock()

	return value
}

func getSessionTokenIdFromCookie(r *WebContext) string {
	c, found := r.Cookie(SessionCookieId(), "")
	if found {
		r.SetCookie(SessionCookieId(), c.Value, time.Hour*24*30)
		return c.Value
	}

	tokenId := toolkit.GenerateRandomString("", 32)
	r.SetCookie(SessionCookieId(), tokenId, time.Hour*24*30)

	return tokenId
}

func (r *WebContext) Session(key string, defs ...interface{}) interface{} {
	tokenId := getSessionTokenIdFromCookie(r)
	var def interface{}
	if len(defs) > 0 {
		def = defs[0]
	}
	return sessions.Get(tokenId, key, def)
}

func (r *WebContext) SetSession(key string, value interface{}) {
	tokenId := getSessionTokenIdFromCookie(r)
	sessions.Set(tokenId, key, value)
}
