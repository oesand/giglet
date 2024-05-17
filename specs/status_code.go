package specs

import (
	"encoding/binary"
	"giglet/safe"
)

type StatusCode struct {
	nocopy safe.NoCopy
	
	code uint16
	text []byte
}

func (status *StatusCode) Append(buffer []byte) []byte {
	buffer = binary.BigEndian.AppendUint16(buffer, status.code)
	buffer = append(buffer, ' ')
	buffer = append(buffer, status.text...)
	return buffer
}

var (
	StatusCodeContinue           = &StatusCode{code: 100, text: []byte("Continue")}
	StatusCodeSwitchingProtocols = &StatusCode{code: 101, text: []byte("Switching Protocols")}
	StatusCodeProcessing         = &StatusCode{code: 102, text: []byte("Processing")}
	StatusCodeEarlyHints         = &StatusCode{code: 103, text: []byte("Early Hints")}

	StatusCodeOK                   = &StatusCode{code: 200, text: []byte("OK")}
	StatusCodeCreated              = &StatusCode{code: 201, text: []byte("Created")}
	StatusCodeAccepted             = &StatusCode{code: 202, text: []byte("Accepted")}
	StatusCodeNonAuthoritativeInfo = &StatusCode{code: 203, text: []byte("Non-Authoritative Information")}
	StatusCodeNoContent            = &StatusCode{code: 204, text: []byte("No Content")}
	StatusCodeResetContent         = &StatusCode{code: 205, text: []byte("Reset Content")}
	StatusCodePartialContent       = &StatusCode{code: 206, text: []byte("Partial Content")}
	StatusCodeMultiStatus          = &StatusCode{code: 207, text: []byte("Multi-Status")}
	StatusCodeAlreadyReported      = &StatusCode{code: 208, text: []byte("Already Reported")}
	StatusCodeIMUsed               = &StatusCode{code: 226, text: []byte("I'm Used")}

	StatusCodeMultipleChoices   = &StatusCode{code: 300, text: []byte("Multiple Choices")}
	StatusCodeMovedPermanently  = &StatusCode{code: 301, text: []byte("Moved Permanently")}
	StatusCodeFound             = &StatusCode{code: 302, text: []byte("Found")}
	StatusCodeSeeOther          = &StatusCode{code: 303, text: []byte("See Other")}
	StatusCodeNotModified       = &StatusCode{code: 304, text: []byte("Not Modified")}
	StatusCodeUseProxy          = &StatusCode{code: 305, text: []byte("Use Proxy")}
	StatusCodeTemporaryRedirect = &StatusCode{code: 307, text: []byte("Temporary Redirect")}
	StatusCodePermanentRedirect = &StatusCode{code: 308, text: []byte("Permanent Redirect")}

	StatusCodeBadRequest                   = &StatusCode{code: 400, text: []byte("Bad Request")}
	StatusCodeUnauthorized                 = &StatusCode{code: 401, text: []byte("Unauthorized")}
	StatusCodePaymentRequired              = &StatusCode{code: 402, text: []byte("Payment Required")}
	StatusCodeForbidden                    = &StatusCode{code: 403, text: []byte("Forbidden")}
	StatusCodeNotFound                     = &StatusCode{code: 404, text: []byte("Not Found")}
	StatusCodeMethodNotAllowed             = &StatusCode{code: 405, text: []byte("Method Not Allowed")}
	StatusCodeNotAcceptable                = &StatusCode{code: 406, text: []byte("Not Acceptable")}
	StatusCodeProxyAuthRequired            = &StatusCode{code: 407, text: []byte("Proxy Auth Required")}
	StatusCodeRequestTimeout               = &StatusCode{code: 408, text: []byte("Request Timeout")}
	StatusCodeConflict                     = &StatusCode{code: 409, text: []byte("Conflict")}
	StatusCodeGone                         = &StatusCode{code: 410, text: []byte("Gone")}
	StatusCodeLengthRequired               = &StatusCode{code: 411, text: []byte("Length Required")}
	StatusCodePreconditionFailed           = &StatusCode{code: 412, text: []byte("Precondition Failed")}
	StatusCodeRequestEntityTooLarge        = &StatusCode{code: 413, text: []byte("Request Entity Too Large")}
	StatusCodeRequestURITooLong            = &StatusCode{code: 414, text: []byte("Request URI Too Long")}
	StatusCodeUnsupportedMediaType         = &StatusCode{code: 415, text: []byte("Unsupported Media Type")}
	StatusCodeRequestedRangeNotSatisfiable = &StatusCode{code: 416, text: []byte("Requested Range Not Satisfiable")}
	StatusCodeExpectationFailed            = &StatusCode{code: 417, text: []byte("ExpectationFailed")}
	StatusCodeTeapot                       = &StatusCode{code: 418, text: []byte("I'm a teapot")}
	StatusCodeMisdirectedRequest           = &StatusCode{code: 421, text: []byte("Misdirected Request")}
	StatusCodeUnprocessableEntity          = &StatusCode{code: 422, text: []byte("Unprocessable Entity")}
	StatusCodeLocked                       = &StatusCode{code: 423, text: []byte("Locked")}
	StatusCodeFailedDependency             = &StatusCode{code: 424, text: []byte("Failed Dependency")}
	StatusCodeTooEarly                     = &StatusCode{code: 425, text: []byte("Too Early")}
	StatusCodeUpgradeRequired              = &StatusCode{code: 426, text: []byte("Upgrade Required")}
	StatusCodePreconditionRequired         = &StatusCode{code: 428, text: []byte("Precondition Required")}
	StatusCodeTooManyRequests              = &StatusCode{code: 429, text: []byte("Too Many Requests")}
	StatusCodeRequestHeaderFieldsTooLarge  = &StatusCode{code: 431, text: []byte("Request Header Fields Too Large")}
	StatusCodeUnavailableForLegalReasons   = &StatusCode{code: 451, text: []byte("Unavailable For Legal Reasons")}

	StatusCodeInternalServerError           = &StatusCode{code: 500, text: []byte("Internal Server Error")}
	StatusCodeNotImplemented                = &StatusCode{code: 501, text: []byte("Not Implemented")}
	StatusCodeBadGateway                    = &StatusCode{code: 502, text: []byte("Bad Gateway")}
	StatusCodeServiceUnavailable            = &StatusCode{code: 503, text: []byte("Service Unavailable")}
	StatusCodeGatewayTimeout                = &StatusCode{code: 504, text: []byte("Gateway Timeout")}
	StatusCodeHTTPVersionNotSupported       = &StatusCode{code: 505, text: []byte("HTTP Version Not Supported")}
	StatusCodeVariantAlsoNegotiates         = &StatusCode{code: 506, text: []byte("Variant Also Negotiates")}
	StatusCodeInsufficientStorage           = &StatusCode{code: 507, text: []byte("Insufficient Storage")}
	StatusCodeLoopDetected                  = &StatusCode{code: 508, text: []byte("Loop Detected")}
	StatusCodeNotExtended                   = &StatusCode{code: 510, text: []byte("Not Extended")}
	StatusCodeNetworkAuthenticationRequired = &StatusCode{code: 511, text: []byte("Network Authentication Required")}
)
