package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context , loan Loan) (Loan, error)
	GetAllLoans(ctx context.Context, status string, order string) ([]Loan, error)
	GetLoanByID(ctx context.Context , loanid string) (Loan , error)
	UpdateLoan(ctx context.Context , id string , loan Loan) (Loan , error)
	DeleteLoan(ctx context.Context , loanid string) error
}

type LoanUseCase interface {
	CreateLoan(c context.Context  , loan LoanRequest , id string) interface{}
	GetAllLoans(c context.Context , id string   , status string  , order string) interface{}
	GetSpecificLoan(ctx context.Context , userid string , loanid string) interface{}
	UpdateLoanStatus(ctx context.Context , userid string , loanid string  , request UpdateLoanRequest) interface{}
	DeleteLoan(ctx context.Context, userid string , loanid string) interface{}
}

type Loan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Userid    string             `json:"userid"`
	Amount    string                `json:"amount"`
	Condition string             `json:"condition"`
}
