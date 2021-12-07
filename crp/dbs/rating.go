package dbs

import "forum/models"

// Rate inserts like/dislike into the table.
func Rate(r *models.Rate) error {
	_, err := conn.Exec(
		`INSERT INTO rate 
		(rate_type, obj_type, uid, obj_id)
		VALUES (?,?,?,?)`,
		r.Type, r.ObjectType, r.UID, r.ObjectID,
	)
	return err
}
