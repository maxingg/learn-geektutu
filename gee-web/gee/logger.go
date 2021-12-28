package gee

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(context *Context) {
		t := time.Now()
		// proess event
		context.Next()
		log.Printf("[%d] %s in %v", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
