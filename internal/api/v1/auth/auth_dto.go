package auth

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firsName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}
