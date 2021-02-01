package response

type UserResponse struct {
	LoggedUserId string `json:"logged_user_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
}
