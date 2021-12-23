package gee

import (
	"fmt"
	"net/http"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandleFunc
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprintf(r.Method, "-", r.URL.Path)
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Printf("404 NOT FOUND: %s\n", r.URL.Path)
	}
}

func (engine *Engine) Get(pattern string, handleFunc HandleFunc) {
	engine.addRoute("GET", pattern, handleFunc)
}

func (engine *Engine) Post(pattern string, handleFunc HandleFunc) {
	engine.addRoute("POST", pattern, handleFunc)
}

// 这里实现了代码复用
func (engine *Engine) addRoute(method string, pattern string, handleFunc HandleFunc) {
	key := fmt.Sprintf(method, "-", pattern)
	engine.router[key] = handleFunc
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandleFunc),
	}
}
