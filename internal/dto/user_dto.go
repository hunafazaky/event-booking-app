package dto

// UserResponse is the public shape of a user, returned anywhere a user
// appears in the API — standalone or nested inside an event/booking.
// Password and internal timestamps are intentionally excluded.
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SignInResponse is returned by POST /signin — the JWT plus the
// authenticated user's public profile in one payload.
type SignInResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
