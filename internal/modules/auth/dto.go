package auth

// This file holds transport (input/output) DTOs.
//
// DTOs decouple the HTTP surface from internal domain models so they can
// evolve independently. Add struct tags matching your validator tags here.

// LoginRequest is the input for a login attempt.
//
// TODO(auth): add validation tags, e.g. `validate:"required,email"` and a
// password `validate:"required,min=8"` once real rules are defined.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the output of a successful login.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	User        User   `json:"user"`
}

// RegisterRequest is the input for registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse is the output of a successful registration.
type RegisterResponse struct {
	User User `json:"user"`
}
