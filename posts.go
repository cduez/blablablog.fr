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
	"path"
	"bytes"
)

type Post struct {
	Slug string
	Name string
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

func ImgUrl(name string) func(string)string {
	return func(filename string) string {
		return generateImgUrl("/images", name, filename)
	}
}

func ImgThumbUrl(name string) func(string)string {
	return func(filename string) string {
		return generateImgUrl("/images", name, "thumbs/" + filename)
	}
}

func generateImgUrl(prefix, name, filename string) string {
	return prefix + "/" + name + "/" + filename
}

func ParseDate(date string) time.Time {
	time, _ := time.Parse("2-1-2006", date);
	return time
}

func NewPosts() Posts {
	posts := Posts{}

	files, _ := filepath.Glob("./posts/*")
	for _, file := range files {
		data, _ := ioutil.ReadFile(file)

		lines := strings.Split(string(data), "\n")
		title := string(lines[0])
		slug  := slugify.Slugify(title)
		name  := strings.TrimSuffix(path.Base(file), filepath.Ext(file))
		date := ParseDate(string(lines[1]))
		country := string(lines[2])
		content := strings.Join(lines[4:len(lines)], "\n")

		funcMap := template.FuncMap{
			"img": ImgUrl(name),
			"imgThumb": ImgThumbUrl(name),
		}

		for k,v := range HelpersFuncs {
			funcMap[k] = v
		}

		var buf bytes.Buffer
		tmpl := template.Must(template.New("post").Funcs(funcMap).Parse(string(content)))
		err := tmpl.Execute(&buf, nil)
		if err != nil {
			panic(err)
		}

		htmlContent := template.HTML(string(blackfriday.MarkdownCommon(buf.Bytes())))

		posts = append(posts, Post{slug, name, title, date, country, htmlContent})
	}

	sort.Sort(posts)

	return posts
}
