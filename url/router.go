package url

import (
	"slices"
	"sort"
	"strings"
)

type UrlRouter[ST any] struct {
	routes map[int][]*Route[ST]
}

func NewRouter[ST any]() *UrlRouter[ST] {
	return &UrlRouter[ST]{
		routes: map[int][]*Route[ST]{},
	}
}

func (router *UrlRouter[ST]) Register(pattern string, stored *ST) {
	route := createRoute(pattern, stored)
	depth := strings.Count(pattern, "/") - 1

	store, exists := router.routes[depth]
	if !exists {
		store = []*Route[ST]{}
	}

	store = append(store, route)

	sort.Slice(store, func(i, n int) bool {
		return len(store[i].names) < len(store[n].names)
	})
	router.routes[depth] = store
}

func (router *UrlRouter[ST]) Include(target *UrlRouter[ST], paths ...string) {
	var prefix string
	if paths != nil {
		prefix = strings.Join(paths, "/")
		if prefix != "" {
			if prefix[0] != '/' {
				panic("prefix \"" + prefix + "\" should not starts with \"/\"")
			} else if prefix[len(prefix) - 1] == '/' {
				panic("prefix \"" + prefix + "\" cannot ends with \"/\"")
			}
		}
	}

	for depth, routes := range target.routes {
		store, exists := router.routes[depth]
		if !exists {
			store = routes
		} else {
			store = slices.Concat(store, routes)
		}
		router.routes[depth] = store

		sort.Slice(store, func(i, n int) bool {
			return len(store[i].names) < len(store[n].names)
		})
	}
}

func (router *UrlRouter[ST]) LookUp(path string) *Route[ST] {
	depth := strings.Count(path, "/") - 1

	store, exists := router.routes[depth]
	if exists {
		for _, route := range store {
			if route.IsMatch(path) {
				return route
			}
		}
	}
	return nil
}

func (router *UrlRouter[ST]) All() []*Route[ST] {
	if len(router.routes) == 0 {
		return nil
	}
	var output []*Route[ST]
	for _, routes := range router.routes {
		if output == nil {
			output = routes;
		} else {
			output = slices.Concat(output, routes)
		}
	}

	return output
}
