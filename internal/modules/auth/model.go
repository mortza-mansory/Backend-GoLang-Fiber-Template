package auth

// User is the auth domain model.
//
// In a real implementation this would map to a users table. It is kept here
// only to illustrate where domain types live.
type User struct {
	ID       string
	Email    string
	Password string // hashed; never plaintext
	Roles    []string
}

// NewUser constructs a User. Placeholder only.
func NewUser(email, passwordHash string) *User {
	return &User{Email: email, Password: passwordHash}
}
