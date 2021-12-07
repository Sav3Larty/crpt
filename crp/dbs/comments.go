package dbs

import "forum/models"

// CreateComment ...
func CreateComment(c *models.Comment) error {
	_, err := conn.Exec(
		`INSERT INTO comment 
		(post_id, uid, text, creation_date)
		VALUES(?,?,?,?)`,
		c.PostID, c.UID, c.Text, c.CreationDate,
	)
	return err
}

// GetComments ...
func GetComments(uid, postID int) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)

	rows, err := conn.Query(
		`SELECT comment.id, comment.post_id, comment.uid, comment.text, comment.creation_date, user.username,
			IFNULL (rate.rate_type,0),
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
				WHERE obj_type = 2 AND obj_id = comment.id
			)
		FROM comment 
		INNER JOIN user ON comment.uid=user.id 
		LEFT JOIN rate ON rate.obj_id = comment.id AND rate.obj_type = 2 AND rate.uid = ? 
		WHERE post_id = ?
		ORDER BY comment.creation_date`,
		uid, postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		comm := &models.Comment{}

		err = rows.Scan(&comm.ID, &comm.PostID, &comm.UID, &comm.Text, &comm.CreationDate, &comm.Author, &comm.UserRate, &comm.Rating)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// DeleteComment ...
func DeleteComment(commentID int) error {
	_, err := conn.Exec(`DELETE FROM comment WHERE id=?`, commentID)
	return err
}
