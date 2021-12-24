package gee

import (
	"fmt"
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandleFunc
}

func (r *router) addRoute(method string, pattern string, handleFunc HandleFunc) {
	log.Printf("Route %4s - %s\n", method, pattern)
	key := fmt.Sprintf(method, "-", pattern)
	r.handlers[key] = handleFunc
}

func (r *router) handle(context *Context) {
	key := fmt.Sprintf(context.Method, "-", context.Path)
	if handler, ok := r.handlers[key]; ok {
		handler(context)
	} else {
		context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandleFunc),
	}
}
