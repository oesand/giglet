package specs

import (
	"errors"
	"io"
	"strconv"
)

type StatusCode struct {
	Code uint16
	Detail []byte
}

func (status *StatusCode) ShouldHaveBody() bool {
	return !(100 <= status.Code && status.Code <= 199 || status.Code == 204 || status.Code == 304)
}

func (status *StatusCode) IsValid() bool {
	return 100 <= status.Code && status.Code < 600
}

func (status *StatusCode) AppendBytes(buffer []byte) []byte {
	buffer = strconv.AppendUint(buffer, uint64(status.Code), 10)
	buffer = append(buffer, ' ')
	buffer = append(buffer, status.Detail...)
	return buffer
}

func (status *StatusCode) WriteAsHeadlineTo(writer io.Writer, is11 bool) error {
	if writer == nil { 
		return errors.New("invalid writer")
	}
	var line []byte
	if is11 {
		line = append(line, httpV11...)
	} else {
		line = append(line, httpV10...)
	}

	line = append(line, ' ')
	line = status.AppendBytes(line)
	line = append(line, directCrlf...)

	_, err := writer.Write(line)
	return err
}

var (
	StatusCodeContinue           = StatusCode{Code: 100, Detail: []byte("Continue")}
	StatusCodeSwitchingProtocols = StatusCode{Code: 101, Detail: []byte("Switching Protocols")}
	StatusCodeProcessing         = StatusCode{Code: 102, Detail: []byte("Processing")}
	StatusCodeEarlyHints         = StatusCode{Code: 103, Detail: []byte("Early Hints")}

	StatusCodeOK                   = StatusCode{Code: 200, Detail: []byte("OK")}
	StatusCodeCreated              = StatusCode{Code: 201, Detail: []byte("Created")}
	StatusCodeAccepted             = StatusCode{Code: 202, Detail: []byte("Accepted")}
	StatusCodeNonAuthoritativeInfo = StatusCode{Code: 203, Detail: []byte("Non-Authoritative Information")}
	StatusCodeNoContent            = StatusCode{Code: 204, Detail: []byte("No Content")}
	StatusCodeResetContent         = StatusCode{Code: 205, Detail: []byte("Reset Content")}
	StatusCodePartialContent       = StatusCode{Code: 206, Detail: []byte("Partial Content")}
	StatusCodeMultiStatus          = StatusCode{Code: 207, Detail: []byte("Multi-Status")}
	StatusCodeAlreadyReported      = StatusCode{Code: 208, Detail: []byte("Already Reported")}
	StatusCodeIMUsed               = StatusCode{Code: 226, Detail: []byte("I'm Used")}

	StatusCodeMultipleChoices   = StatusCode{Code: 300, Detail: []byte("Multiple Choices")}
	StatusCodeMovedPermanently  = StatusCode{Code: 301, Detail: []byte("Moved Permanently")}
	StatusCodeFound             = StatusCode{Code: 302, Detail: []byte("Found")}
	StatusCodeSeeOther          = StatusCode{Code: 303, Detail: []byte("See Other")}
	StatusCodeNotModified       = StatusCode{Code: 304, Detail: []byte("Not Modified")}
	StatusCodeUseProxy          = StatusCode{Code: 305, Detail: []byte("Use Proxy")}
	StatusCodeTemporaryRedirect = StatusCode{Code: 307, Detail: []byte("Temporary Redirect")}
	StatusCodePermanentRedirect = StatusCode{Code: 308, Detail: []byte("Permanent Redirect")}

	StatusCodeBadRequest                   = StatusCode{Code: 400, Detail: []byte("Bad Request")}
	StatusCodeUnauthorized                 = StatusCode{Code: 401, Detail: []byte("Unauthorized")}
	StatusCodePaymentRequired              = StatusCode{Code: 402, Detail: []byte("Payment Required")}
	StatusCodeForbidden                    = StatusCode{Code: 403, Detail: []byte("Forbidden")}
	StatusCodeNotFound                     = StatusCode{Code: 404, Detail: []byte("Not Found")}
	StatusCodeMethodNotAllowed             = StatusCode{Code: 405, Detail: []byte("Method Not Allowed")}
	StatusCodeNotAcceptable                = StatusCode{Code: 406, Detail: []byte("Not Acceptable")}
	StatusCodeProxyAuthRequired            = StatusCode{Code: 407, Detail: []byte("Proxy Auth Required")}
	StatusCodeRequestTimeout               = StatusCode{Code: 408, Detail: []byte("Request Timeout")}
	StatusCodeConflict                     = StatusCode{Code: 409, Detail: []byte("Conflict")}
	StatusCodeGone                         = StatusCode{Code: 410, Detail: []byte("Gone")}
	StatusCodeLengthRequired               = StatusCode{Code: 411, Detail: []byte("Length Required")}
	StatusCodePreconditionFailed           = StatusCode{Code: 412, Detail: []byte("Precondition Failed")}
	StatusCodeRequestEntityTooLarge        = StatusCode{Code: 413, Detail: []byte("Request Entity Too Large")}
	StatusCodeRequestURITooLong            = StatusCode{Code: 414, Detail: []byte("Request URI Too Long")}
	StatusCodeUnsupportedMediaType         = StatusCode{Code: 415, Detail: []byte("Unsupported Media Type")}
	StatusCodeRequestedRangeNotSatisfiable = StatusCode{Code: 416, Detail: []byte("Requested Range Not Satisfiable")}
	StatusCodeExpectationFailed            = StatusCode{Code: 417, Detail: []byte("ExpectationFailed")}
	StatusCodeTeapot                       = StatusCode{Code: 418, Detail: []byte("I'm a teapot")}
	StatusCodeMisdirectedRequest           = StatusCode{Code: 421, Detail: []byte("Misdirected Request")}
	StatusCodeUnprocessableEntity          = StatusCode{Code: 422, Detail: []byte("Unprocessable Entity")}
	StatusCodeLocked                       = StatusCode{Code: 423, Detail: []byte("Locked")}
	StatusCodeFailedDependency             = StatusCode{Code: 424, Detail: []byte("Failed Dependency")}
	StatusCodeTooEarly                     = StatusCode{Code: 425, Detail: []byte("Too Early")}
	StatusCodeUpgradeRequired              = StatusCode{Code: 426, Detail: []byte("Upgrade Required")}
	StatusCodePreconditionRequired         = StatusCode{Code: 428, Detail: []byte("Precondition Required")}
	StatusCodeTooManyRequests              = StatusCode{Code: 429, Detail: []byte("Too Many Requests")}
	StatusCodeRequestHeaderFieldsTooLarge  = StatusCode{Code: 431, Detail: []byte("Request Header Fields Too Large")}
	StatusCodeUnavailableForLegalReasons   = StatusCode{Code: 451, Detail: []byte("Unavailable For Legal Reasons")}

	StatusCodeInternalServerError           = StatusCode{Code: 500, Detail: []byte("Internal Server Error")}
	StatusCodeNotImplemented                = StatusCode{Code: 501, Detail: []byte("Not Implemented")}
	StatusCodeBadGateway                    = StatusCode{Code: 502, Detail: []byte("Bad Gateway")}
	StatusCodeServiceUnavailable            = StatusCode{Code: 503, Detail: []byte("Service Unavailable")}
	StatusCodeGatewayTimeout                = StatusCode{Code: 504, Detail: []byte("Gateway Timeout")}
	StatusCodeHTTPVersionNotSupported       = StatusCode{Code: 505, Detail: []byte("HTTP Version Not Supported")}
	StatusCodeVariantAlsoNegotiates         = StatusCode{Code: 506, Detail: []byte("Variant Also Negotiates")}
	StatusCodeInsufficientStorage           = StatusCode{Code: 507, Detail: []byte("Insufficient Storage")}
	StatusCodeLoopDetected                  = StatusCode{Code: 508, Detail: []byte("Loop Detected")}
	StatusCodeNotExtended                   = StatusCode{Code: 510, Detail: []byte("Not Extended")}
	StatusCodeNetworkAuthenticationRequired = StatusCode{Code: 511, Detail: []byte("Network Authentication Required")}
)
