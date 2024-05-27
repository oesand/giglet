package specs

import (
	"net/url"
)

type Query url.Values

func ParseQuery(query string) (Query, error) {
	obj, err := url.ParseQuery(query)
	return Query(obj), err
}
