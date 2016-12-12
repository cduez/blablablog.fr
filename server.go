package main

import (
	"io"
	"net/http"
	"html/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "base", "")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name + ".tmpl", data)
}

func main() {
	e := echo.New()
	e.ShutdownTimeout = 3

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} -> ${status} in ${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("./public/*.tmpl")),
	}

	e.Static("/", "assets")
	e.Renderer = t

	e.GET("/", Index)

	e.Logger.Fatal(e.Start(":1323"))
}
