package gee

import (
	"net/http"
)

type HandleFunc func(*Context)

type Engine struct {
	router *router
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := newContext(w, r)
	engine.router.handle(context)
}

func (engine *Engine) GET(pattern string, handleFunc HandleFunc) {
	engine.addRoute("GET", pattern, handleFunc)
}

func (engine *Engine) POST(pattern string, handleFunc HandleFunc) {
	engine.addRoute("POST", pattern, handleFunc)
}

// 这里实现了代码复用
func (engine *Engine) addRoute(method string, pattern string, handleFunc HandleFunc) {
	engine.router.addRoute(method, pattern, handleFunc)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}
