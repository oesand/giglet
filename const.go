package giglet

import "errors"

type RequestHandler func(request *HttpRequest) *HttpResponse

const (
	defaultLineLengthLimit = 1024
)

var (
	errorTooLarge = errors.New("too large")
)