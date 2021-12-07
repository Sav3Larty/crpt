package internal

import (
	"html/template"
)

//GetTempl returns templates.
func GetTempl() (*template.Template, error) {
	templ, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		return nil, err
	}
	return templ, nil
}
