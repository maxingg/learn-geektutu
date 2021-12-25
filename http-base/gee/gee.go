package gee

import (
	"log"
	"net/http"
)

type HandleFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandleFunc
		parent      *RouterGroup
		engine      *Engine
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := newContext(w, r)
	engine.router.handle(context)
}

func (group *RouterGroup) GET(pattern string, handleFunc HandleFunc) {
	group.addRoute("GET", pattern, handleFunc)
}

func (group *RouterGroup) POST(pattern string, handleFunc HandleFunc) {
	group.addRoute("POST", pattern, handleFunc)
}

func (group *RouterGroup) addRoute(method string, comp string, handleFunc HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handleFunc)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
