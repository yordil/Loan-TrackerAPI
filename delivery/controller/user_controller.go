package controller

import (
	"loantracker/domain"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUseCase domain.UserUseCase
}


func (u *UserController) RegisterUser(c *gin.Context) {
	var user domain.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, domain.ErrorResponse{Message: "Invalid request", Status: 400})
		return
	}

	resp := u.UserUseCase.CreateUser(c, user)
	HandleResponse(c , resp)
	
	
}

func (u* UserController) VerifyUser(c *gin.Context) {

	token := c.Query("token")

	if token == ""{
		response := domain.ErrorResponse {
			Message: "Invalid Token", 
			Status: 400,
		}

		HandleResponse(c , response)
	}

	response := u.UserUseCase.VeryfyUser(c , token)

	HandleResponse(c , response)

}

func (u *UserController) Login(c *gin.Context) {
	var user domain.LoginRequest
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, domain.ErrorResponse{Message: "Invalid request", Status: 400})
		return
	}

	resp := u.UserUseCase.Login(c, user)
	HandleResponse(c , resp)
}


func (u *UserController) ForgotPassword(c *gin.Context) {
	var userEmail domain.ForgotPasswordRequest

	token := c.Query("token")

	if token == "" {

		err := c.ShouldBindJSON(&userEmail)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		response := u.UserUseCase.ForgotPassword(c, userEmail)
		HandleResponse(c, response)
		return
	} else {
		var password domain.ResetPasswordRequest
		err := c.ShouldBindJSON(&password)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		response := u.UserUseCase.ResetPassword(c, password, token)

		HandleResponse(c, response)
	}
}


//  refresh Token endPoint

func (u *UserController) RefreshToken(c *gin.Context) {

	

	var refreshToken domain.RefreshTokenRequest

	err := c.ShouldBindJSON(&refreshToken)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response := u.UserUseCase.RefreshToken(c , refreshToken)
	
	HandleResponse(c, response)

}


//  profile mydata 

func (u *UserController) Profile(c *gin.Context) {
	id := c.GetString("user_id")
	
	response := u.UserUseCase.Profile(c, id)
	HandleResponse(c, response)
}


func (u *UserController) GetAllUsers(c *gin.Context) {

	id := c.GetString("user_id")
	response := u.UserUseCase.GetAllUsers(c , id)
	HandleResponse(c, response)
}

func (u *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	adminid := c.GetString("user_id")
	response := u.UserUseCase.DeleteUser(c, id , adminid)
	HandleResponse(c, response)
}