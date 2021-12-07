package dbs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // using driver for sql
)

// conn is holding pool of connetions for database
var conn *sql.DB

// NewConnect ...
func NewConnect() error {
	db, err := sql.Open("sqlite3", "file:forum.s3db?_auth&_auth_user=Dawrld&_auth_pass=Alibi&_auth_crypt=sha256&_foreign_keys=on")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	tables := []string{
		`CREATE TABLE  IF NOT EXISTS "user" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"username"	TEXT UNIQUE,
			"email"	TEXT UNIQUE,
			"password"	TEXT
		)`,

		`CREATE TABLE  IF NOT EXISTS "session" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"uid"	INTEGER,
			"uuid"	TEXT UNIQUE,
			"status"	INTEGER DEFAULT 0,
			"datetime"	DATETIME,
			FOREIGN KEY ("uid") REFERENCES user ("id") ON DELETE CASCADE
		)`,

		`CREATE TABLE  IF NOT EXISTS "obj_type" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"name"	TEXT UNIQUE
		)`,

		`CREATE TABLE  IF NOT EXISTS "post" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"uid"	INTEGER,
			"text"	TEXT,
			"image" TEXT,
			"creation_date"	DATETIME,
			FOREIGN KEY ("uid") REFERENCES user ("id") ON DELETE CASCADE
		)`,

		`CREATE TABLE	IF NOT EXISTS "category" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"name"	TEXT UNIQUE
		)`,

		`CREATE TABLE  IF NOT EXISTS "post_category" (
			"id"	INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			"post_id"	INTEGER,
			"category_id"	INTEGER,
			FOREIGN KEY ("post_id") REFERENCES post ("id") ON DELETE CASCADE
			FOREIGN KEY ("category_id") REFERENCES category ("id") ON DELETE CASCADE
			UNIQUE("post_id", "category_id")
		)`,

		`CREATE TRIGGER IF NOT EXISTS "session_update"
			BEFORE INSERT
			ON session
			BEGIN
				UPDATE session SET status = 0 WHERE uid = NEW.uid;
			END;
		`,
	}

	for _, v := range tables {
		_, err = db.Exec(v)
		if err != nil {
			return err
		}
	}

	fmt.Println("Connected to the database")
	conn = db

	go func() {
		for {
			if err := CleanSessions(); err != nil {
				log.Println(err)
			}
			time.Sleep(10 * time.Minute)
		}
	}()
	return nil
}
