package models

type SignupRequest struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Id string `json:"id"`
}

type LoginResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Response struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    LoginResponse `json:"data"`
}
