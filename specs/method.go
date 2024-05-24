package specs

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodOptions HttpMethod = "OPTIONS"
	HttpMethodHead    HttpMethod = "HEAD"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodTrace   HttpMethod = "TRACE"
)

func (method HttpMethod) IsValid() bool {
	return method == HttpMethodGet ||
		method == HttpMethodPost ||
		method == HttpMethodPut ||
		method == HttpMethodDelete ||
		method == HttpMethodOptions ||
		method == HttpMethodHead ||
		method == HttpMethodPatch ||
		method == HttpMethodTrace
}

func (method HttpMethod) IsPostable() bool {
	return method == HttpMethodPost ||
		method == HttpMethodPut ||
		method == HttpMethodDelete ||
		method == HttpMethodPatch
}
