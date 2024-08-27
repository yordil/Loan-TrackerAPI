package domain

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefrshTokenResponse struct {
	AcessToken string `json:"access_token"`
	Status     int    `json:"status"`
}

type ProfileResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	User    User   `json:"user"`
}

type AllUserResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Users   []User `json:"users"`
}

type LoanResponse struct {
	Status int `json:"status"`
	Data   Loan
}

type AllLoanResponse struct {
	Status int    `json:"status"`
	Loans  []Loan `json:"loans"`
}

type CreateLoanResponse struct {
	Status    int    `json:"status"`
	Condition string `json:"condition"`
}