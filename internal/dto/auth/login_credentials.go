package authdto

type LoginCredentials struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}
