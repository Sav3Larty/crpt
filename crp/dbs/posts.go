package dbs

import (
	"database/sql"
	"forum/models"
)

// PostType is representation for filter
type PostType int64

const (
	SimplePost PostType = iota
	FavouritePost
	MyPost
)

// CreatePost using SQL transactions for insertion in two tables.
// In case of any error, it rolls back any change.
func CreatePost(p *models.Post) error {

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO post (uid, text, creation_date, image) VALUES(?,?,?,?)",
		p.UID, p.Text, p.CreationDate, p.Image,
	)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	p.ID = int(id)

	stmt, err := tx.Prepare("INSERT INTO post_category (post_id, category_id) VALUES(?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, cat := range p.CategoryID {
		_, err := stmt.Exec(id, cat)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetRecentPosts collects data from multiple tables and returns it as formatted posts.
func GetRecentPosts(uid, amount int, postType PostType) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, amount)
	query :=
		`SELECT	post.id, post.uid, post.text, post.creation_date, post.image, user.username, 
		IFNULL (rate.rate_type,0),
		(
			SELECT	COUNT (*)
			FROM comment
			WHERE post.id = comment.post_id
		),
		(
			SELECT IFNULL
				(
					SUM (CASE WHEN rate_type = 1 THEN 1
							WHEN rate_type = 2 THEN -1
							ELSE 0
						END
						),
				0)
			FROM rate
			WHERE obj_type = 1 AND obj_id = post.id
		)
		FROM post
		INNER JOIN user ON post.uid = user.id
		`

	var (
		rows *sql.Rows
		err  error
	)

	if postType == MyPost {
		rows, err = conn.Query(
			query+
				` AND post.uid=?
				LEFT JOIN rate ON rate.obj_id = post.id AND rate.obj_type = 1 
				ORDER BY post.creation_date DESC
				LIMIT ?`,
			uid, amount,
		)
		if err != nil {
			return nil, err
		}

	} else {
		rows, err = conn.Query(
			query+
				`LEFT JOIN rate ON rate.obj_id = post.id AND rate.obj_type = 1 AND rate.uid = ? 
				ORDER BY post.creation_date DESC
				LIMIT ?`,
			uid, amount,
		)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.ID, &post.UID, &post.Text, &post.CreationDate, &post.Image, &post.Author, &post.UserRate, &post.Comments, &post.Rating)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, p := range posts {
		err := setPostCategories(p)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func setPostCategories(p *models.Post) error {
	rows, err := conn.Query(
		`SELECT post_category.category_id, category.name 
		FROM post_category
		INNER JOIN category ON category.id=post_category.category_id 
		WHERE post_category.post_id=?`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			catID   int
			catName string
		)

		err = rows.Scan(&catID, &catName)
		if err != nil {
			return err
		}

		p.CategoryID = append(p.CategoryID, catID)
		p.CategoryName = append(p.CategoryName, catName)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// GetPost finds a post by user id and post id.
func GetPost(uid, id int) (*models.Post, error) {
	post := &models.Post{}
	if err := conn.QueryRow(
		`SELECT	post.id, post.uid, post.text, post.creation_date, post.image, user.username FROM post
		INNER JOIN user ON post.uid = user.id
		WHERE post.id = ?`,
		id,
	).Scan(
		&post.ID, &post.UID, &post.Text, &post.CreationDate, &post.Image, &post.Author,
	); err != nil {
		return nil, err
	}

	if err := setPostCategories(post); err != nil {
		return post, err
	}

	return post, nil
}

//GetFavouritePost from database where user put like.
func GetFavouritePost(uid int) ([]*models.Post, error) {
	rows, err := conn.Query(
		`SELECT obj_id FROM rate WHERE rate_type =1 AND obj_type=1 AND uid=?`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getPostsByParam(uid, rows)
}

//GetPostByCategory filter all posts by categories.
func GetPostByCategory(categoryID, uid int) ([]*models.Post, error) {
	rows, err := conn.Query(
		`SELECT post_id FROM post_category WHERE category_id =?`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getPostsByParam(uid, rows)
}

func getPostsByParam(uid int, rows *sql.Rows) ([]*models.Post, error) {
	postIDs := make([]int, 0)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	for _, postID := range postIDs {
		post, err := GetPost(uid, int(postID))
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// DeletePost is causing cascade deletion of post and related data.
func DeletePost(postID int) error {
	_, err := conn.Exec(`DELETE FROM post WHERE id=?`, postID)
	return err
}
