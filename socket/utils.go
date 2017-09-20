package socket

import (
	"fmt"
	"html/template"
	"net/http"
)

//
//
//
func render(w http.ResponseWriter, contentName string, context Context) {

	context.Static = STATIC_URL

	templates := []string{"templates/base.html"}

	contentFilename := fmt.Sprintf("templates/%s.html", contentName)

	templates = append(templates, contentFilename)

	t, err := template.ParseFiles(templates...)
	if err != nil {
		fmt.Print("template parsing error: ", err)
	}

	if err = t.Execute(w, context); err != nil {
		fmt.Print("template executing error: ", err)
	}
}
