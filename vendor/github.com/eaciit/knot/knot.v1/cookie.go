package knot

import (
	"net/http"
	"net/url"
	"time"
)

var DefaultCookieExpire time.Duration

func (r *WebContext) initCookies() {
	if r.cookies == nil {
		r.cookies = make(map[string]*http.Cookie)
	}
}

func (r *WebContext) Cookie(name string, def string) (*http.Cookie, bool) {
	r.initCookies()

	// first search on new cookies
	c, exist := r.cookies[name]

	// when not found, try to search on request cookies
	if exist == false {
		var err error
		c, err = r.Request.Cookie(name)
		if err == nil {
			exist = true
		}
	}

	// when not exist and default is set
	// put cookie with default expire time
	if exist == false && len(def) > 0 {
		if int(DefaultCookieExpire) == 0 {
			DefaultCookieExpire = 30 * 24 * time.Hour
		}
		r.SetCookie(name, def, DefaultCookieExpire)
	}

	return c, exist
}

func (r *WebContext) SetCookie(name string, value string, expiresAfter time.Duration) *http.Cookie {
	r.initCookies()

	c := &http.Cookie{}
	c.Name = name
	c.Value = value
	c.Path = "/"
	u, e := url.Parse(r.Request.URL.String())
	if e == nil {
		c.Expires = time.Now().Add(expiresAfter)
		c.Domain = u.Host
	}

	r.cookies[name] = c

	return c
}

func (r *WebContext) Cookies() map[string]*http.Cookie {
	r.initCookies()

	return r.cookies
}
