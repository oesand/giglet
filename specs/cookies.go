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
	builder := bytes.Buffer{}
	
	builder.WriteString(cookie.Name)
	builder.WriteByte('=')
	builder.WriteString(cookie.Value)

	if cookie.MaxAge > 0 {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeyMaxAge)
		builder.WriteByte('=')
		builder.Write(strconv.AppendUint(nil, cookie.MaxAge, 10))
	} else if !cookie.Expires.IsZero() {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeyExpires)
		builder.WriteByte('=')
		builder.Write(cookie.Expires.UTC().AppendFormat(nil, TimeFormat))
	}

	if len(cookie.Domain) > 0 {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeyDomain)
		builder.WriteByte('=')
		builder.WriteString(cookie.Domain)
	}

	if len(cookie.Path) > 0 {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeyPath)
		builder.WriteByte('=')
		builder.WriteString(cookie.Path)
	}

	if cookie.HttpOnly {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeyHTTPOnly)
	}

	if cookie.Secure {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeySecure)
	}

	if len(cookie.SameSite) > 0 {
		builder.Write(cookieDelimiter)
		builder.Write(cookieKeySameSite)
		builder.WriteByte('=')
		builder.WriteString(string(cookie.SameSite))
	}

	return builder.Bytes()
}
