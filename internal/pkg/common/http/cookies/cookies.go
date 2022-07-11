package cookies

import (
	"github.com/chmike/securecookie"
)

type Cookies struct {
	SecureCookie    map[string]*securecookie.Obj
	SecureCookieKey []byte
}

var App = new(Cookies)
var CookieSettings = make(map[string]securecookie.Params)

func SetupCookies(list map[string]securecookie.Params) error {
	App.SecureCookie = make(map[string]*securecookie.Obj)
	App.SecureCookieKey = securecookie.MustGenerateRandomKey()
	CookieSettings = list
	for k, v := range CookieSettings {
		sc, err := securecookie.New(k, App.SecureCookieKey, v)
		if err != nil {
			return err
		}
		App.SecureCookie[k] = sc
	}
	return nil
}

func ReplenishExpiredCookie(key string, t int) *securecookie.Obj {
	p := CookieSettings[key]
	p.MaxAge = t
	return securecookie.MustNew(key, App.SecureCookieKey, p)
}
