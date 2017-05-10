package main

import (
	"html/template"
	"bytes"
)

var HelpersFuncs = template.FuncMap{
	"youtube": Youtube,
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
