package main

import (
	"fmt"
	"gee"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	//r := gee.New()
	//r.Use(gee.Logger()) // global midlleware
	//r.SetFuncMap(template.FuncMap{
	//	"FormatAsDate": FormatAsDate,
	//})
	//r.LoadHTMGlob("templates/*")
	//r.Static("/assets", "./static")
	//
	//stu1 := &student{Name: "Geektutu", Age: 20}
	//stu2 := &student{Name: "Jack", Age: 22}
	//
	//r.GET("/", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "css.tmpl.tmpl", nil)
	//})
	//
	//r.GET("/students", func(context *gee.Context) {
	//	context.HTML(http.StatusOK, "arr.tmpl", gee.H{
	//		"title":  "gee",
	//		"stuArr": [2]*student{stu1, stu2},
	//	})
	//})
	//
	//r.GET("/date", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "custom_func.tmpl.tmpl", gee.H{
	//		"title": "gee",
	//		"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
	//	})
	//})

	r := gee.Default()
	r.GET("/", func(context *gee.Context) {
		context.String(http.StatusOK, "Hello Geektutu\n")
	})

	r.GET("/panic", func(context *gee.Context) {
		names := []string{"geektutu"}
		context.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}
