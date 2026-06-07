package response

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
}
