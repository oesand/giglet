package url

func ParseQuery(query string) Query {
	return Query{}
}

type Query map[string][]string

func (query Query) Get(key string) string {
	arr := query[key]
	if len(arr) == 0 {
		return ""
	}
	return arr[0]
}

func (query Query) Set(key, value string) {
	query[key] = []string{value}
}

func (query Query) Add(key, value string) {
	query[key] = append(query[key], value)
}

func (query Query) Delete(key string) {
	delete(query, key)
}

func (query Query) Has(key string) bool {
	_, ok := query[key]
	return ok
}
