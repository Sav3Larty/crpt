package models

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Comment is absract representation of comment.
// Including author details and related data.
type Comment struct {
	ID           int
	PostID       int
	UID          int
	Author       string
	Text         string
	Rating       int
	UserRate     int
	CreationDate time.Time
}

// Validate is checking size of comment.
// Valid comment should contain letter besides spaces.
func (c *Comment) Validate() error {
	if strings.TrimSpace(c.Text) == "" {
		return fmt.Errorf("empty comment")
	}
	if utf8.RuneCountInString(c.Text) >= 300 {
		return fmt.Errorf("comment is too long")
	}
	return nil
}
