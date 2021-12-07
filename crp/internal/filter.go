package internal

import (
	"forum/dbs"
	"forum/models"
	"strconv"
)

func filter(filterParam string, user *models.User) ([]*models.Post, error) {
	var posts []*models.Post
	categoryID, err := strconv.Atoi(filterParam)
	if err == nil {
		if exists, err := dbs.ExistsCategory(categoryID); err != nil || !exists {
			return nil, errCategoryNotExists
		}
		posts, err = dbs.GetPostByCategory(categoryID, user.ID)
		if err != nil {
			return nil, err
		}
	} else if filterParam == "favourites" {
		if !user.Authorized {
			return nil, errUserNotAuth
		}
		posts, err = dbs.GetFavouritePost(user.ID)
		if err != nil {
			return nil, err
		}
	} else if filterParam == "my-posts" {
		if !user.Authorized {
			return nil, errUserNotAuth
		}
		posts, err = dbs.GetRecentPosts(user.ID, 1000, dbs.MyPost)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errCategoryNotExists
	}
	return posts, nil
}
