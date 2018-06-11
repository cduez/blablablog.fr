package main

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
)

func HelpersFuncs(name string) template.FuncMap {
	return template.FuncMap{
		"youtube":  Youtube,
		"picture":  Picture(name),
		"pictures": Pictures(name),
	}
}

func Youtube(videoId string) template.HTML {
  var output bytes.Buffer
	tmpl := template.Must(template.New("youtube").ParseFiles("public/helpers/youtube.tmpl"))

	err := tmpl.Execute(&output, struct{ Id string }{videoId})
	if err != nil {
		panic(err)
	}

	return template.HTML(output.String())
}

type PictureParams struct {
	Name string
	Id   string
}

func Picture(name string) func(string) template.HTML {
	return func(imageId string) template.HTML {
		var output bytes.Buffer
		tmpl := template.Must(template.New("picture").ParseFiles("public/helpers/picture.tmpl"))

		err := tmpl.Execute(&output, PictureParams{name, imageId})
		if err != nil {
			panic(err)
		}

		return template.HTML(output.String())
	}
}

func Pictures(name string) func(string, string) template.HTML {
	// "1-15,6,23-24" "caption"
	return func(groups, caption string) template.HTML {
		var buffer template.HTML = "<figure>"
		splittedGroup := strings.Split(groups, ",")

		for _, rg := range splittedGroup {
			var first, last int

			if strings.Contains(rg, "-") {
				splitted := strings.Split(rg, "-")
				first, _ = strconv.Atoi(splitted[0])
				last, _ = strconv.Atoi(splitted[1])
			} else {
				first, _ = strconv.Atoi(rg)
				last = first
			}

			for i := first; i <= last; i++ {
				buffer += Picture(name)(strconv.Itoa(i))
			}
		}

		if caption != "" {
			buffer += template.HTML("<figcaption>" + caption + "</figcaption>")
		}
		buffer += "</figure>"

		return buffer
	}
}
