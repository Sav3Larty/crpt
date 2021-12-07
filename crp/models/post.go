package models

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Post is absract representation of post.
// Including author details and related data.
type Post struct {
	ID           int
	CategoryID   []int
	CategoryName []string
	UID          int
	Author       string
	Text         string
	Image        string
	Rating       int
	Comments     int
	UserRate     int
	CreationDate time.Time
}

// Validate returns non nil error in case of empty post or when post is too long.
func (p *Post) Validate() error {
	if strings.TrimSpace(p.Text) == "" {
		return fmt.Errorf("empty post")
	}
	if utf8.RuneCountInString(p.Text) >= 500 {
		return fmt.Errorf("post is too long")
	}
	if len(p.CategoryID) == 0 {
		return fmt.Errorf("no categories set")
	}
	return nil
}
