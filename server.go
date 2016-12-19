package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sort"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/russross/blackfriday"
	"github.com/Machiel/slugify"
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

type CustomRenderer struct {}

func FormatDate(date time.Time) string {
	return fmt.Sprintf("%d %s %d", date.Day(), translateMonth[date.Month().String()], date.Year())
}

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

func ParseDate(date string) time.Time {
	time, _ := time.Parse("2-1-2006", date);
	return time
}

type Post struct {
	Slug string
	Title string
	Date time.Time
	Country string // ISO 3166
	Content template.HTML
}
type Posts []Post

func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Posts) Less(i, j int) bool { return (p[i].Date).After(p[j].Date) }

func Index(c echo.Context) error {
	posts := getPosts()
	return c.Render(http.StatusOK, "index", posts)
}

func ViewPost(c echo.Context) error {
	postName := c.Param("slug")
	posts := getPosts()
	post := findPostBySlug(posts, postName)

	return c.Render(http.StatusOK, "post", post)
}

func getPosts() []Post {
	posts := Posts{}

	files, _ := filepath.Glob("./posts/*")
	for _, file := range files {
		data, _ := ioutil.ReadFile(file)

		lines := strings.Split(string(data), "\n")
		title := string(lines[0])
		slug  := slugify.Slugify(title)
		date := ParseDate(string(lines[1]))
		country := string(lines[2])
		content := strings.Join(lines[4:len(lines)], "\n")
		htmlContent := template.HTML(string(blackfriday.MarkdownCommon([]byte(content))))
		posts = append(posts, Post{slug, title, date, country, htmlContent})
	}

	sort.Sort(posts)

	return posts
}

func findPostBySlug(posts []Post, slug string) Post {
	for _, post := range posts {
		if slug == post.Slug {
			return post
		}
	}
	return Post{}
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
