package dbs

import (
	"forum/models"
	"time"
)

// FillDatabase with example data
func FillDatabase() error {
	users := [][]string{
		{
			"alibi@mail.ru",
			"alibi",
			"Alibi123",
		},

		{
			"damir@gmail.com",
			"Dawrld",
			"D4mir",
		},
	}

	for _, v := range users {
		user := &models.User{
			Email:    v[0],
			Username: v[1],
			Password: v[2],
		}
		err := CreateUser(user)
		if err != nil {
			return err
		}
	}

	objects := []string{
		"post",
		"comment",
	}

	for _, obj := range objects {
		_, err := conn.Exec("INSERT INTO obj_type (name) VALUES(?)", obj)
		if err != nil {
			return err
		}
	}

	rates := []string{
		"like",
		"dislike",
	}

	for _, rate := range rates {
		_, err := conn.Exec("INSERT INTO rate_type (name) VALUES(?)", rate)
		if err != nil {
			return err
		}
	}

	cats := []string{
		"anime",
		"games",
		"memes",
		"music",
		"sport",
	}

	for _, v := range cats {
		err := CreateCategory(v)
		if err != nil {
			return err
		}
	}

	texts := []string{
		"lorem ipsum",
		"Dossan mal",
		"Curiosity killed the cat",
		"Поехали",
	}

	for _, v := range texts {
		post := &models.Post{
			CategoryID:   []int{1, 2, 3},
			UID:          1,
			Text:         v,
			CreationDate: time.Now(),
		}
		err := CreatePost(post)
		if err != nil {
			return err
		}

	}
	return nil
}
