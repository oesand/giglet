package specs

import (
	"io"
	"mime/multipart"
)

type SM map[string]string
type TM map[string]any
type AR []any

type Stream io.Reader
type MultipartForm multipart.Form
type UploadFile multipart.FileHeader


const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

var (
	cookieParamDelimiter    = []byte("; ")
	cookieKeyExpires        = []byte("Expires")
	cookieKeyDomain         = []byte("Domain")
	cookieKeyPath           = []byte("Path")
	cookieKeyHTTPOnly       = []byte("HttpOnly")
	cookieKeySecure         = []byte("Secure")
	cookieKeyMaxAge         = []byte("Max-Age")
	cookieKeySameSite       = []byte("SameSite")
)
