package models

import "regexp"

//User is absraction over user.
type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Authorized bool
}

// Validate is verifying user credentials for registration.
// Such as username/password length and email format.
func (u *User) Validate() bool {
	if len(u.Username) < 5 || len(u.Password) < 5 || len(u.Email) < 5 {
		return false
	}

	if re := regexp.MustCompile(`^([a-z0-9_-]+\.)*[a-z0-9_-]+@[a-z0-9_-]+(\.[a-z0-9_-]+)*\.[a-z]{2,6}$`); !re.MatchString(u.Email) {
		return false
	}
	return true
}
