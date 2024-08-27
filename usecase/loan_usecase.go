package usecase

import (
	"context"
	"loantracker/domain"
	"loantracker/infrastructure"
	"time"
)

type LoanUseCase struct {
	LoanRepository  domain.LoanRepository
	contextTimeout  time.Duration
	config          *infrastructure.Config
	userRepo 	  domain.UserRepository
}

func NewLoanUsecase(ur domain.LoanRepository, timeout time.Duration, config *infrastructure.Config, userrepo domain.UserRepository) domain.LoanUseCase {
	return &LoanUseCase{
		LoanRepository:  ur,
		contextTimeout:  timeout,
		config: config,
		userRepo: userrepo,
	}
}

func (l *LoanUseCase) CreateLoan(c context.Context , loanamount domain.LoanRequest , id string) interface{} {

	ctx, cancel := context.WithTimeout(c, l.contextTimeout)
	defer cancel()

	loan := domain.Loan{}

	loan.Amount = loanamount.Amount

	loan.Userid = id
	loan.Condition = "pending"
	loan, err := l.LoanRepository.CreateLoan(ctx, loan)
	if err != nil {
		return domain.ErrorResponse{Message: "Could not create loan", Status: 400}
	}

	return domain.CreateLoanResponse{Status: 200, Condition: loan.Condition}
	
}


//  get all loans

func (l *LoanUseCase) GetAllLoans(c context.Context , id string , status , order string ) interface{} {
	ctx, cancel := context.WithTimeout(c, l.contextTimeout)
	defer cancel()

	// check if user is admin or not by find user by id and check the role
	user, err := l.userRepo.FindUserByID(ctx , id)
	

	if err != nil {
		return domain.ErrorResponse{Message: "Could not get user", Status: 400}
	}

	if user.Role != "admin" {
		return domain.ErrorResponse{Message: "You are not authorized to view this page", Status: 400}
	}

	loans, err := l.LoanRepository.GetAllLoans(ctx , status , order)
	if err != nil {
		return domain.ErrorResponse{Message: err.Error(), Status: 400}
	}

	return domain.AllLoanResponse{Status: 200, Loans: loans}
}


func (l *LoanUseCase) GetSpecificLoan(ctx context.Context , userid string , loanid string) interface {} {

	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	// check if user is admin or not by find user by id and check the role
	user, err := l.userRepo.FindUserByID(ctx , userid)

	if err != nil {
		return domain.ErrorResponse{Message: "Could not get user", Status: 400}
	}

	if user.Role != "admin" {
		return domain.ErrorResponse{Message: "You are not authorized to view this page", Status: 400}
	}

	// CALL GET BY ID 

	loan , err := l.LoanRepository.GetLoanByID(ctx , loanid)

	if err != nil {
		return domain.ErrorResponse{Message: "loan Not Found" , Status: 404}
	}

	return domain.LoanResponse{Status : 200 , Data :loan}
	
}


func (l *LoanUseCase) UpdateLoanStatus(ctx context.Context , userid string , loanid string , request domain.UpdateLoanRequest) interface{} {

	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	// check if user is admin or not by find user by id and check the role
	user, err := l.userRepo.FindUserByID(ctx , userid)

	if err != nil {
		return domain.ErrorResponse{Message: "Could not get user", Status: 400}
	}

	if user.Role != "admin" {
		return domain.ErrorResponse{Message: "You are not authorized to Approve or reject the loan", Status: 400}
	}

	loan , err := l.LoanRepository.GetLoanByID(ctx , loanid)

	if err != nil {
		return domain.ErrorResponse{Message: "Loan Not Found" , Status: 404}
	}

	
	if request.Action == "approve" {
		loan.Condition = "approved"

		//  call update in here 

		loan , err = l.LoanRepository.UpdateLoan(ctx , loanid , loan)

		if err != nil {
			return domain.ErrorResponse{Message: "Could not approve loan", Status: 400}
		}
		
		return domain.SuccessResponse{Message: "Loan approved Successfully" , Status: 200}

	} else if request.Action == "reject" {
		loan.Condition = "rejected"

		//  call update in here

		loan , err = l.LoanRepository.UpdateLoan(ctx , loanid , loan)

		if err != nil {
			return domain.ErrorResponse{Message: "Could not reject loan", Status: 400}
		}

		return domain.SuccessResponse{Message: "Loan rejected Successfully" , Status: 200}
	}

	return domain.ErrorResponse{Message: "Invalid action you can either reject or approve loan" , Status: 200}
	

}


// delete loan 

func (l *LoanUseCase) DeleteLoan(ctx context.Context, userid string , loanid string) interface{} {

	ctx , cancel := context.WithTimeout(ctx, l.contextTimeout)

	defer cancel()

	// check if user is admin or not by find user by id and check the role
	user, err := l.userRepo.FindUserByID(ctx , userid)

	if err != nil {
		return domain.ErrorResponse{Message: "Could not get user", Status: 400}
	}

	if user.Role != "admin" { 
		return domain.ErrorResponse{Message: "You are not authorized to delete the loan", Status: 400}
	}

	_ , err = l.LoanRepository.GetLoanByID(ctx , loanid)

	if err != nil {
		return domain.ErrorResponse{Message: "Loan Not Found" , Status: 404}
	}

	err = l.LoanRepository.DeleteLoan(ctx , loanid)

	if err != nil {
		return domain.ErrorResponse{Message: "Could not delete loan", Status: 400}
	}

	return domain.SuccessResponse{Message: "Loan deleted Successfully" , Status: 200}

}