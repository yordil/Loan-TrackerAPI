package domain


type RegisterUserRequest struct { 
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsVerified bool `json:"is_verified"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type RefreshTokenRequest struct { 
	RefreshToken string `json:"refresh_token"`
}
type LoginRequest struct { 
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoanRequest struct { 
	Amount string `json:"amount"`
}

type UpdateLoanRequest struct { 
	Action string `json:"action"`
}