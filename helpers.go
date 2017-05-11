package main

import (
	"html/template"
	"bytes"
	"strconv"
)

func HelpersFuncs(name string) template.FuncMap{
	return template.FuncMap{
		"youtube": Youtube,
		"picture": Picture(name),
		"groupPicture": GroupPicture(name),
	}
}

func Youtube(videoId string) template.HTML {
	var output bytes.Buffer
	tmpl := template.Must(template.New("youtube").ParseFiles("public/helpers/youtube.tmpl"))

	err := tmpl.Execute(&output, struct{Id string}{videoId,})
	if err != nil {
		panic(err)
	}

	return template.HTML(output.String())
}

type PictureParams struct {
	Name string
	Id string
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

func GroupPicture(name string) func(int, int, string) template.HTML {
	var buffer template.HTML = "<figure>"

	return func(first, last int, caption string) template.HTML {
		for i := first; i <= last; i++ {
			buffer += Picture(name)(strconv.Itoa(i))
		}
		if caption != "" {
			buffer += template.HTML("<figcaption>" + caption + "</figcaption>")
		}
		buffer += "</figure>"

		return buffer
	}
}
