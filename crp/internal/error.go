package internal

import (
	"errors"
	"html/template"
	"net/http"
)

var (
	errCategoryNotExists = errors.New("this category not exists")
	errUserNotAuth       = errors.New("user not authorised")
)

//Error handle for errors, setting header and executing message to template.
func Error(w http.ResponseWriter, templ *template.Template, errorStatusCode int, msg string) {
	w.WriteHeader(errorStatusCode)
	if err := templ.ExecuteTemplate(w, "error", struct{ ErrorMessage string }{msg}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
