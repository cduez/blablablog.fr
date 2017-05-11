package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
	"os"
	"io/ioutil"
	"strings"
	"crypto/sha1"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/joho/godotenv"
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

func IsCurrentPage(currentPage string) func(string) bool {
	return func(page string) bool {
		return page == currentPage
	}
}

func ContainerCurrentPage(currentPage string) string {
	switch currentPage {
	case "map":
		return "container container--large"
	}

	return "container"
}

func styleSHA1() string {
	fmt.Println("copute")
	f, err := os.Open("assets/stylesheets/style.css")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

type CustomRenderer struct {}

func (t *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	funcMap := template.FuncMap{
		"format": FormatDate,
		"isCurrentPage": IsCurrentPage(name),
		"containerCurrentPage": func() string { return ContainerCurrentPage(name) },
		"styleSHA1": styleSHA1,
	}

	tmpl := template.New("base.tmpl").Funcs(funcMap)

	tmplParsed := template.Must(tmpl.ParseFiles("public/base.tmpl", "public/common.tmpl" , "public/" + name + ".tmpl"))
	return tmplParsed.Execute(w, data)
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

func getPoints() (points [][]string) {
	data, _ := ioutil.ReadFile("./points")
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		coord := strings.Split(line, ":")
		if len(coord) > 1 {
			points = append(points, coord[0:2])
		}
	}

	return
}

func ViewMap(c echo.Context) error {
	data :=  map[string]interface{}{
		"token": os.Getenv("MAPBOX_TOKEN"),
		"points": getPoints(),
	}
	return c.Render(http.StatusOK, "map", data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.ShutdownTimeout = 3

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} -> ${status} in ${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	e.Static("/", "assets")
	e.Renderer = &CustomRenderer{}

	e.GET("/posts/:slug", ViewPost)
	e.GET("/map", ViewMap)
	e.GET("/", Index)

	e.Logger.Fatal(e.Start(":1666"))
}
