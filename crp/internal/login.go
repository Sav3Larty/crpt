package internal

import (
	"forum/dbs"
	"forum/models"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Login handler for logging.
func Login(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postLogin(w, r, templ)
		case http.MethodGet:
			getLogin(w, r, templ)
		default:
			Error(w, templ, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}
	})
}

func getLogin(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	user := r.Context().Value(ctxKeyUser).(*models.User)

	if user.Authorized {
		err := dbs.DeactivateSession(user.ID)
		if err != nil {
			Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}

	if err := templ.ExecuteTemplate(w, "login", nil); err != nil {
		Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func postLogin(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	password := ""
	username := ""
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	query, err := url.ParseQuery(string(bytes))
	if err != nil {
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	for i, v := range query {
		switch i {
		case "username":
			username = v[0]
		case "password":
			password = v[0]
		default:
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
	}

	user, err := dbs.FindUser(username, password)
	if err != nil {
		Error(w, templ, http.StatusUnauthorized, "Incorrect username or password")
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
}
