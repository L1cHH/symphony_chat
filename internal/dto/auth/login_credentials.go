package authdto

type LoginCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
