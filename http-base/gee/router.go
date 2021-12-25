package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handleFunc HandleFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handleFunc
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if ok {
		n := root.search(searchParts, 0)
		if n != nil {
			parts := parsePattern(n.pattern)
			for index, part := range parts {
				if part[0] == ':' {
					params[part[1:]] = searchParts[index]
				}
				if part[0] == '*' && len(part) > 1 {
					params[part[1:]] = strings.Join(searchParts[index:], "/")
				}
			}
			return n, params
		}
	}
	return nil, nil
}

func (r *router) handle(context *Context) {
	n, params := r.getRoute(context.Method, context.Path)
	if n != nil {
		context.Params = params
		key := context.Method + "-" + n.pattern
		r.handlers[key](context)
	} else {
		context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandleFunc),
		roots:    make(map[string]*node),
	}
}
