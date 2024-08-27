package controller

import (
	"fmt"
	"loantracker/domain"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	LoanUseCase domain.LoanUseCase
}


//  create loan conttoller

func (u *LoanController) CreateLoan(c * gin.Context){

	userid := c.GetString("user_id")

	var loan domain.LoanRequest

	err := c.ShouldBindJSON(&loan)

	if err != nil {
		c.JSON(400, domain.ErrorResponse{Message: err.Error(), Status: 400 })
		return
	}

	response := u.LoanUseCase.CreateLoan(c, loan, userid)

	HandleResponse(c , response)
}

func (u *LoanController) GetAllLoans(c *gin.Context) {

	status := c.DefaultQuery("status", "all")
	order := c.DefaultQuery("order", "asc")
	id := c.GetString("user_id")
	fmt.Println(id  , "******************************************")
	response := u.LoanUseCase.GetAllLoans(c , id , status , order)

	HandleResponse(c , response) 

}

func (u *LoanController) GetSpecificLoan(c *gin.Context) {
	loanid := c.Param("id")
	userid := c.GetString("user_id")
	response := u.LoanUseCase.GetSpecificLoan(c , userid , loanid)

	HandleResponse(c , response)
}


func (u *LoanController) UpdateLoanStatus(c *gin.Context) {
	loanid := c.Param("id")
	userid := c.GetString("user_id")
	var updateRequest domain.UpdateLoanRequest 
	err := c.ShouldBindJSON(&updateRequest)

	if err != nil {
		c.JSON(400, domain.ErrorResponse{Message: err.Error(), Status: 400 })
		return
	}



	response := u.LoanUseCase.UpdateLoanStatus(c , userid , loanid , updateRequest)

	HandleResponse(c , response)
}

//  delete loan

func (u *LoanController) DeleteLoan(c *gin.Context) {
	loanid := c.Param("id")
	userid := c.GetString("user_id")

	response := u.LoanUseCase.DeleteLoan(c , userid , loanid)

	HandleResponse(c , response)
}