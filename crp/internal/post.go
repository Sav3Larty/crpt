package internal

import (
	"errors"
	"fmt"
	"forum/dbs"
	"forum/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 20 MB
const maxSize = 20 * 1024 * 1024

var (
	errLarge      = errors.New("File size is too large")
	errType       = errors.New("File content type is not supported")
	errInternal   = errors.New("Internal server error")
	errBadRequest = errors.New("Bad request")
)

// Post handler for posts.
func Post(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addPost(w, r, templ)
		case http.MethodGet:
			openPost(w, r, templ)
		default:
			Error(w, templ, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}
	})
}

func openPost(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	if r.URL.Path == "/post/" {
		Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	i, err := strconv.Atoi(r.URL.Path[6:])
	if i <= 0 || err != nil {
		Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	user := r.Context().Value(ctxKeyUser).(*models.User)

	post, err := dbs.GetPost(user.ID, i)
	fmt.Println(post, err)
	if err != nil {
		Error(w, templ, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	page := struct {
		User     *models.User
		Post     *models.Post
		Comments []*models.Comment
	}{User: user, Post: post}

	if err := templ.ExecuteTemplate(w, "post", page); err != nil {
		Error(w, templ, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

}

func addPost(w http.ResponseWriter, r *http.Request, templ *template.Template) {

	user := r.Context().Value(ctxKeyUser).(*models.User)

	if !user.Authorized {
		Error(w, templ, http.StatusUnauthorized, "please login to post")
		return
	}

	var (
		text       string
		categories []int
	)

	if r.Form == nil {
		r.ParseMultipartForm(32 << 20)
	}

	for k, v := range r.Form {
		switch k {
		case "category":
			for _, cat := range v {
				n, err := strconv.Atoi(cat)
				if err != nil {
					Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
					return
				}
				categories = append(categories, n)
			}
		case "text":
			text = v[0]

		case "Image":
			continue

		default:
			Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
	}

	post := &models.Post{
		CategoryID:   categories,
		UID:          user.ID,
		Text:         text,
		CreationDate: time.Now(),
	}

	if status, err := imageUpload(r, post); err != nil {
		Error(w, templ, status, err.Error())
		return
	}

	if err := post.Validate(); err != nil {
		Error(w, templ, http.StatusBadRequest, err.Error())
		return
	}

	if err := dbs.CreatePost(post); err != nil {
		fmt.Println(err)
		Error(w, templ, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	link := fmt.Sprintf("/post/%v", post.ID)
	http.Redirect(w, r, link, http.StatusFound)
}

func imageUpload(r *http.Request, post *models.Post) (int, error) {
	file, handler, err := r.FormFile("Image")

	if err == http.ErrMissingFile {
		return http.StatusOK, nil
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer file.Close()

	if handler.Size > maxSize {
		return http.StatusRequestEntityTooLarge, errLarge
	}

	contentTypes := handler.Header["Content-Type"]
	if len(contentTypes) != 1 || (contentTypes[0] != "image/png" && contentTypes[0] != "image/jpeg" && contentTypes[0] != "image/gif") {
		return http.StatusBadRequest, errType
	}

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + handler.Filename

	dst, errFileCreate := os.Create("static/media/" + fileName)
	if errFileCreate != nil {
		return http.StatusInternalServerError, errInternal
	}
	defer dst.Close()

	if _, errFileCopy := io.Copy(dst, file); errFileCopy != nil {
		return http.StatusInternalServerError, errInternal
	}

	f, err := os.Open("static/media/" + fileName)
	if err != nil {
		return http.StatusInternalServerError, errInternal
	}
	defer f.Close()

	contentType, err := FileContentType(f)
	if err != nil {
		return http.StatusBadRequest, errBadRequest
	}

	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" {
		err := os.Remove("static/media/" + fileName)
		if err != nil {
			fmt.Println(err)
		}
		return http.StatusInternalServerError, errInternal
	}

	post.Image = "/static/media/" + fileName
	return http.StatusOK, nil
}

func FileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
