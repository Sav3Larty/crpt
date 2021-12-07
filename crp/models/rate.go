package models

// Rate is an absraction over like/dislike.
type Rate struct {
	ID         int // rate id number
	Type       int // 1 - like, 2 - dislike
	ObjectType int // 1 - post, 2 - comment
	UID        int // user id
	ObjectID   int // reference id to object
}
