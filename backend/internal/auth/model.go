package auth

// GoogleUser matches the JSON response from Google's UserInfo API
type GoogleUser struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
