package main

import (
	"time"
	"html/template"
	"io/ioutil"
	"strings"
	"github.com/russross/blackfriday"
	"github.com/Machiel/slugify"
	"path/filepath"
	"sort"
)

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

func (p Posts) FindBySlug(slug string) Post {
	for _, post := range p {
		if slug == post.Slug {
			return post
		}
	}
	return Post{}
}

func NewPosts() Posts {
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
