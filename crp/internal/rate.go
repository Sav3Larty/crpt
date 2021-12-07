package internal

import (
	"fmt"
	"forum/dbs"
	"forum/models"
	"html/template"
	"net/http"
	"strconv"
)

// Rating is handler for like/dislike actions.
func Rating(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		if r.Method != http.MethodPost {
			Error(w, templ, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		user := r.Context().Value(ctxKeyUser).(*models.User)

		if !user.Authorized {
			Error(w, templ, http.StatusUnauthorized, "please login to like")
			return
		}

		err := r.ParseForm()
		if err != nil {
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		var (
			rateTypeID   int
			objectTypeID int
			objectID     int
		)

		for k, v := range r.PostForm {
			switch k {
			case "action":
				rateTypeID, err = strconv.Atoi(v[0])

			case "objType":
				objectTypeID, err = strconv.Atoi(v[0])

			case "objID":
				objectID, err = strconv.Atoi(v[0])
			default:
				Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}

			if err != nil {
				Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}
		}

		action := &models.Rate{
			Type:       rateTypeID,
			ObjectType: objectTypeID,
			UID:        user.ID,
			ObjectID:   objectID,
		}

		if err := dbs.Rate(action); err != nil {
			fmt.Println(err)
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		link := fmt.Sprintf("%v#%v", r.Header.Get("Referer"), objectID)
		http.Redirect(w, r, link, http.StatusFound)

	})
}
