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
		"pictureL": PictureL(name),
		"pictureP": PictureP(name),
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
	Orientation string
}

func PictureFN(name, imageId, orientation string) template.HTML {
	var output bytes.Buffer
	tmpl := template.Must(template.New("picture").ParseFiles("public/helpers/picture.tmpl"))

	err := tmpl.Execute(&output, PictureParams{name, imageId, orientation})
	if err != nil {
		panic(err)
	}

	return template.HTML(output.String())
}

func Picture(name string) func(string, string) template.HTML {
	return func(imageId, orientation string) template.HTML {
		return PictureFN(name, imageId, orientation)
	}
}

func PictureL(name string) func(string) template.HTML {
	return func(imageId string) template.HTML {
		return PictureFN(name, imageId, "landscape")
	}
}

func PictureP(name string) func(string) template.HTML {
	return func(imageId string) template.HTML {
		return PictureFN(name, imageId, "portrait")
	}
}

func GroupPicture(name string) func(int, int) template.HTML {
	var buffer template.HTML

	return func(first, last int) template.HTML {
		for i := first; i <= last; i++ {
			buffer += PictureL(name)(strconv.Itoa(i))
		}

		return buffer
	}
}
