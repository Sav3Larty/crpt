package dbs

import (
	"time"
)

// UserIDBySession is checking active cookies and returns user id.
func UserIDBySession(session string) (int, error) {
	var userID int
	if err := conn.QueryRow(
		"SELECT uid FROM session WHERE status = 1 AND uuid = ?", session,
	).Scan(&userID); err != nil {
		return userID, err
	}
	return userID, nil
}

// CreateSession is ...
func CreateSession(uid int, uuid string, date time.Time) error {
	_, err := conn.Exec(
		"INSERT INTO session (uid, uuid, status, datetime) VALUES (?,?,?,?)",
		uid, uuid, 1, date,
	)
	return err

}

// CleanSessions is designed to be run within goroutine and periodically update status.
func CleanSessions() error {
	time := time.Now()
	_, err := conn.Exec("UPDATE session SET status=0 WHERE status=1 AND datetime < ?", time)
	return err
}

// DeactivateSession ...
func DeactivateSession(uid int) error {
	_, err := conn.Exec("UPDATE session SET status=0 WHERE status=1 AND uid=?", uid)
	return err
}
