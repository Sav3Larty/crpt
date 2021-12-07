package internal

import (
	"fmt"
	"forum/dbs"
	"forum/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Comment handler for adding comments.
func Comment(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		case http.MethodPost:
			addComment(w, r, templ)
		default:
			Error(w, templ, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}
	})
}

func addComment(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	user := r.Context().Value(ctxKeyUser).(*models.User)

	if !user.Authorized {
		Error(w, templ, http.StatusUnauthorized, "please login to comment")
		return
	}

	err := r.ParseForm()
	if err != nil {
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	var (
		text   string
		postID int
	)

	for k, v := range r.PostForm {
		switch k {
		case "postID":
			postID, err = strconv.Atoi(v[0])

		case "text":
			text = v[0]

		default:
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		if err != nil {
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
	}

	comment := &models.Comment{
		PostID:       postID,
		UID:          user.ID,
		Text:         text,
		CreationDate: time.Now(),
	}

	if err := comment.Validate(); err != nil {
		Error(w, templ, http.StatusBadRequest, err.Error())
		return
	}

	if err := dbs.CreateComment(comment); err != nil {
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	link := fmt.Sprintf("../post/%v", postID)
	http.Redirect(w, r, link, http.StatusFound)
}
