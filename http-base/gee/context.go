package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	statusCode int
}

func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) status(code int) {
	c.statusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Data(code int, data []byte) {
	c.status(code)
	c.Writer.Write(data)
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}
