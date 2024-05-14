package giglet

type RequestHandler func(request *HttpRequest) *HttpResponse
