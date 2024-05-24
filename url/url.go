package url

func ParseUrl(url string) (*Url, error) {
	return &Url{}, nil
}

type Url struct {
	scheme, username, password,
	host, port, path, query, hash string

	queryParams Query
}

func (url *Url) Scheme() string {
	return url.scheme
}

func (url *Url) Username() string {
	return url.username
}

func (url *Url) Password() string {
	return url.password
}

func (url *Url) Host() string {
	return url.host
}

func (url *Url) Port() string {
	return url.port
}

func (url *Url) Path() string {
	return url.path
}

func (url *Url) Query() string {
	return url.query
}

func (url *Url) QueryParams() Query {
	if url.queryParams == nil {
		url.queryParams = ParseQuery(url.query)
	}
	return url.queryParams
}

func (url *Url) Hash() string {
	return url.hash
}
