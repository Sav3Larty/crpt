package internal

import (
	"forum/dbs"
	"forum/models"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Registration handler for registration.
func Registration(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postReg(w, r, templ)
		case http.MethodGet:
			getReg(w, r, templ)
		default:
			Error(w, templ, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}
	})
}

func getReg(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	if err := templ.ExecuteTemplate(w, "login", nil); err != nil {
		Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func postReg(w http.ResponseWriter, r *http.Request, templ *template.Template) {

	defer r.Body.Close()
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	query, err := url.ParseQuery(string(bytes))
	if err != nil {
		Error(w, templ, http.StatusBadRequest, err.Error())
		return
	}
	user := &models.User{}
	for i, v := range query {
		switch i {
		case "password":
			user.Password = v[0]
		case "username":
			user.Username = v[0]
		case "email":
			user.Email = v[0]
		default:
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
	}
	if !user.Validate() {
		Error(w, templ, http.StatusBadRequest, "Not correct email, username or password. all of them should have at least 5 char and email should be in correct format")
		return
	}
	if err := dbs.CreateUser(user); err != nil {
		Error(w, templ, http.StatusBadRequest, "User with this username or email exist, please try use another")
		return
	}

	expiration := time.Now().Add(time.Hour)
	u2 := uuid.NewV4()
	cookie := http.Cookie{Name: "session", Value: u2.String(), Expires: expiration}

	err = dbs.CreateSession(user.ID, u2.String(), expiration)
	if err != nil {
		Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)

	if err := templ.ExecuteTemplate(w, "login", nil); err != nil {
		Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

}
