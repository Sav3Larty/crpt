package dbs

import (
	"database/sql"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

// FindUser ...
func FindUser(username, password string) (*models.User, error) {
	user := &models.User{}
	var passwordFromDB string

	if err := conn.QueryRow(
		"SELECT id, username, password FROM user WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &passwordFromDB); err != nil {
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password))
	return user, err
}

// FindUserBySession returns User model with fields: ID, Username, Authorization status.
func FindUserBySession(session string) (*models.User, error) {
	user := &models.User{Authorized: true}

	err := conn.QueryRow(
		`SELECT session.uid, user.username
		FROM session
		INNER JOIN user ON session.uid=user.id
		WHERE session.status = 1 AND session.uuid = ?`, session,
	).Scan(&user.ID, &user.Username)

	if err != nil {
		user.Authorized = false
		user.Username = "Guest"
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return user, err
}

// DeleteUser is causing cascade delete of all data associated to user.
// Including posts, comments, likes.
func DeleteUser(username string) error {
	_, err := conn.Exec(`DELETE FROM user WHERE username=?`, username)
	return err
}

// CreateUser ....
func CreateUser(u *models.User) error {
	hashedpaswrd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	res, err := conn.Exec(
		"INSERT INTO user (username,email,password) VALUES(?,?,?)",
		u.Username, u.Email, string(hashedpaswrd),
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	u.ID = int(id)
	return err
}
