package specs

import (
	"bytes"
	"strconv"
	"time"
)

type CookieSameSite string

const (
	CookieSameSiteLaxMode		CookieSameSite = "Lax"
	CookieSameSiteStrictMode	CookieSameSite = "Strict"
	CookieSameSiteNoneMode		CookieSameSite = "None"
)

type Cookie struct {
	Name    	string
	Value  		string

	Domain 		string
	MaxAge 		uint64
	Expires 	time.Time
	Path   		string

	HttpOnly 	bool
	Secure   	bool
	SameSite 	CookieSameSite
}

func (cookie *Cookie) Bytes() []byte {
	var buf bytes.Buffer
	
	buf.WriteString(cookie.Name)
	buf.WriteByte('=')
	buf.WriteString(cookie.Value)

	if cookie.MaxAge > 0 {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeyMaxAge)
		buf.WriteByte('=')
		buf.Write(strconv.AppendUint(nil, cookie.MaxAge, 10))
	} else if !cookie.Expires.IsZero() {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeyExpires)
		buf.WriteByte('=')
		buf.Write(cookie.Expires.UTC().AppendFormat(nil, TimeFormat))
	}

	if len(cookie.Domain) > 0 {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeyDomain)
		buf.WriteByte('=')
		buf.WriteString(cookie.Domain)
	}

	if len(cookie.Path) > 0 {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeyPath)
		buf.WriteByte('=')
		buf.WriteString(cookie.Path)
	}

	if cookie.HttpOnly {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeyHTTPOnly)
	}

	if cookie.Secure {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeySecure)
	}

	if len(cookie.SameSite) > 0 {
		buf.Write(cookieDelimiter)
		buf.Write(cookieKeySameSite)
		buf.WriteByte('=')
		buf.WriteString(string(cookie.SameSite))
	}

	return buf.Bytes()
}
