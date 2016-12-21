package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var translateMonth map[string]string = map[string]string {
	"January": "Janvier",
	"February": "Février",
	"March": "Mars",
	"April": "Avril",
	"May": "Mai",
	"June": "Juin",
	"July": "Juillet",
	"August": "Août",
	"September": "Septembre",
	"October": "Octobre",
	"November": "Novembre",
	"December": "Décembre",
}

func FormatDate(date time.Time) string {
	return fmt.Sprintf("%d %s %d", date.Day(), translateMonth[date.Month().String()], date.Year())
}

func ParseDate(date string) time.Time {
	time, _ := time.Parse("2-1-2006", date);
	return time
}

type CustomRenderer struct {}

func (t *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	funcMap := template.FuncMap{
		"format": FormatDate,
	}

	tmpl := template.New("base.tmpl").Funcs(funcMap)

	if name == "post" {
		tmplPost := template.Must(tmpl.ParseFiles("public/base.tmpl", "public/post.tmpl"))
		return tmplPost.Execute(w, data)
	} else {
		tmplIndex := template.Must(tmpl.ParseFiles("public/base.tmpl", "public/index.tmpl"))
		return tmplIndex.Execute(w, data)
	}
}

func Index(c echo.Context) error {
	posts := NewPosts()
	return c.Render(http.StatusOK, "index", posts)
}

func ViewPost(c echo.Context) error {
	postName := c.Param("slug")
	posts := NewPosts()
	post := posts.FindBySlug(postName)

	return c.Render(http.StatusOK, "post", post)
}

func main() {
	e := echo.New()
	e.ShutdownTimeout = 3

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} -> ${status} in ${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	e.Static("/", "assets")
	e.Renderer = &CustomRenderer{}

	e.GET("/posts/:slug", ViewPost)
	e.GET("/", Index)

	e.Logger.Fatal(e.Start(":1666"))
}
