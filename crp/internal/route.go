package internal

import (
	"forum/models"
	"html/template"
	"net/http"
)

func Register(templ *template.Template, limiter *models.Limiter) *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	})
	mux.Handle("/home", Middleware(limiter, Index(templ)))
	mux.Handle("/login", Middleware(limiter, Login(templ)))
	mux.Handle("/registration", Middleware(limiter, Registration(templ)))
	mux.Handle("/post/", Middleware(limiter, Post(templ)))
	return mux
}
