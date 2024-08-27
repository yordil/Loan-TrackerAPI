package usecase

import (
	"context"
	"loantracker/domain"
	"loantracker/infrastructure"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct { 
	UserRepository domain.UserRepository;
	contextTimeout time.Duration
	passwordService domain.PasswordService
	config *infrastructure.Config
}

func NewUserUseCase(ur domain.UserRepository, timeout time.Duration , passwordservice domain.PasswordService , config *infrastructure.Config) domain.UserUseCase {
	return &UserUseCase{
		UserRepository: ur,
		contextTimeout: timeout,
		passwordService: passwordservice,
		config: config,
	}
}

func (u *UserUseCase) CreateUser(c context.Context , user domain.User) interface{} { 

	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return  domain.ErrorResponse{Message: "All fields are required" , Status: 400}
	}

	// check if user already exists
	userExists, _ := u.UserRepository.FindUserByEmail(ctx, user.Email)
	if userExists.Email != "" {
		return domain.ErrorResponse{Message: "User already exists" , Status: 400}
	}

	user.IsVerified = false

	// hash password
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return domain.ErrorResponse{Message: "Error hashing password" , Status: 500}
	}
	user.Password = hashedPassword

	// call the generator token 
	token , err := infrastructure.GenerateVerifyToken()

	if err != nil {
		return domain.ErrorResponse{Message: "Error generating token" , Status: 500}
	}
	user.VerificationToken = token
	
	user.VerificationExpiration = time.Now().Add(time.Minute * 10)
	user.ID = primitive.NewObjectID()

	_, err = u.UserRepository.CreateUser(ctx ,  user)
	if err != nil {
		return domain.ErrorResponse{Message: "Error creating user" , Status: 500}
	}

	// send email to user

	err = infrastructure.SendVerifyEmail(user.Email , token)

	if err != nil {
		return domain.ErrorResponse{Message: "Error sending email" , Status: 500}
	}

	return domain.SuccessResponse{Message: "User created successfully verify your account" , Status: 200}

}

func (u *UserUseCase) FindUserByEmail(c context.Context , email string) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	_ , err := u.UserRepository.FindUserByEmail(ctx , email)
	if err != nil { 
		return domain.ErrorResponse{Message: "Error finding user" , Status: 500}
	}

	return domain.SuccessResponse{Message: "User found" , Status: 200}

}


func (u *UserUseCase) VeryfyUser(c context.Context , token string) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	user , err := u.UserRepository.FindUserByVerificationToken(ctx , token)

	
	if err != nil {
		return domain.ErrorResponse{Message: "Error finding user" , Status: 500}
	}

	if user.VerificationExpiration.Before(time.Now()) {
		return domain.ErrorResponse{Message: "Token expired" , Status: 400}
	}

	// update the user
	user.IsVerified = true
	user.ResetPasswordToken = ""
	user.VerificationToken = ""
	user.VerificationExpiration = time.Time{}
	user.ResetPasswordExpires = time.Time{}
	_ , err = u.UserRepository.UpdateUser(ctx , user)
	if err != nil {
		return domain.ErrorResponse{Message: "Error verifying user" , Status: 500}
	}

	return domain.SuccessResponse{Message: "User verified successfully" , Status: 200}

}

func (u *UserUseCase) Login(c context.Context , user domain.LoginRequest) interface{} { 

	ctx , cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()

	exisitingUser := domain.User{}

	exisitingUser , err := u.UserRepository.FindUserByEmail(ctx , user.Email)

	if err != nil { 
		return domain.ErrorResponse{Message: "User Not found" , Status: 500}
	}

	if !exisitingUser.IsVerified { 
		return domain.ErrorResponse{Message: "User not verified" , Status: 400}
	}

	// compare password
	match := u.passwordService.ComparePassword(user.Password , exisitingUser.Password )

	if !match {
		return domain.ErrorResponse{Message: "Invalid password" , Status: 400}
	}

	stringid := exisitingUser.ID.Hex()
	refresh := time.Duration(100) * time.Hour
	access := time.Duration(15) * time.Minute
	acessToken , err := infrastructure.GenerateToken(stringid , access , u.config.AccessTokenSecret)
	if err != nil {
		return domain.ErrorResponse{Message: "Error generating token" , Status: 500}
	}

	refreshToken , err := infrastructure.GenerateToken(stringid , refresh , u.config.RefreshTokenSecret)
	if err != nil {
		return domain.ErrorResponse{Message: "Error generating token" , Status: 500}
	}
	exisitingUser.RefreshToken = refreshToken
	// update
	_ , err = u.UserRepository.UpdateUser(ctx , exisitingUser)

	if err != nil {
		return domain.ErrorResponse{Message: "Error in setting Referesh Token" , Status: 500}
	}

	return domain.LoginResponse{Message: "Logged in Sucessfully" , AccessToken: acessToken , RefreshToken: refreshToken , Status: 200}

}


func (u *UserUseCase) ForgotPassword(c context.Context, email domain.ForgotPasswordRequest) interface{} {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check if user exists
	existing, err := u.UserRepository.FindUserByEmail(ctx, email.Email)
	if err != nil {
		return domain.ErrorResponse{Message: "User not found", Status: 404}
	}

	
	if existing.ResetPasswordToken != "" && time.Now().Before(existing.ResetPasswordExpires) {
		difftime := time.Until(existing.ResetPasswordExpires)
		return domain.ErrorResponse{Message: "Reset token already sent Please wait for " + strconv.FormatFloat(difftime.Minutes(), 'f', -1, 64) + " to resend reset token", Status: 400}
	}

	token, err := infrastructure.GenerateResetToken()

	if err != nil {
		return domain.ErrorResponse{Message: "Error generating reset token", Status: 500}
	}

	

	expiration := time.Now().Add(time.Minute * 50)

	_, err = u.UserRepository.SetResetToken(ctx, email, token, expiration)

	if err != nil {
		return domain.ErrorResponse{Message: "Error saving reset token", Status: 500}
	}

	// send reset email

	err = infrastructure.SendResetEmail(email.Email, token)

	if err != nil {
		return domain.ErrorResponse{Message: "Error sending reset email", Status: 500}
	}

	return domain.SuccessResponse{Message: "Reset email sent", Status: 200}

}

func (u *UserUseCase) ResetPassword(c context.Context, password domain.ResetPasswordRequest, token string) interface{} {

	// check the password validity
	if err := u.passwordService.ValidatePassword(password.Password); err != nil {
		return domain.ErrorResponse{Message: err.Error(), Status: 400}
	}

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check if the ResetToken is Set

	user, err := u.UserRepository.FindUserByResetToken(ctx, token)

	if err != nil {
		return domain.ErrorResponse{Message: "Invalid reset token", Status: 400}
	}

	// check if token is expired
	if user.ResetPasswordExpires.Before(time.Now()) { 
		return domain.ErrorResponse{Message: "Reset token expired", Status: 400}
	}
	


	hashedPassword, err := u.passwordService.HashPassword(password.Password)

	if err != nil {
		return domain.ErrorResponse{Message: "Error hashing password", Status: 500}
	}

	// hash the password

	user.Password = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetPasswordExpires = time.Time{}

	// update user
	_, err = u.UserRepository.UpdateUser(ctx, user)

	if err != nil {
		return &domain.ErrorResponse{Message: "Error in Reseting the Password", Status: 500}
	}

	return &domain.SuccessResponse{Message: "Password Reset Sucessfully", Status: 200}
}

//  refresh Token endPoint

func (u *UserUseCase) RefreshToken(c context.Context,  refreshToken domain.RefreshTokenRequest) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)

	defer cancel()
	
	// check if user exists
	user, err := u.UserRepository.FindUserByRefreshToken(ctx, refreshToken.RefreshToken)
	
	if err != nil { 
		return domain.ErrorResponse{Message: "User Not found" , Status: 500}
	}
	
	// generate new access Token 
	//  check if the refresh token is expired or not or changed 
	if user.RefreshToken == "" {
		return domain.ErrorResponse{Message: "Invalid refresh token" , Status: 400}
	}
	if user.RefreshToken != refreshToken.RefreshToken {
		return domain.ErrorResponse{Message: "Invalid refresh token" , Status: 400}
	}

	// check if the refresh token is expired
	check , _ := infrastructure.IsAuthorized(refreshToken.RefreshToken , u.config.RefreshTokenSecret)

	if !check {
		return domain.ErrorResponse{Message: "Refresh Token Expired" , Status : 400}
	}

	// generate new access token and refresh token
	
	stringid := user.ID.Hex()
	access := time.Duration(15) * time.Minute
	acessToken , err := infrastructure.GenerateToken(stringid , access , u.config.AccessTokenSecret)
	if err != nil {
		return domain.ErrorResponse{Message: "Error generating token" , Status: 500}
	}

	return domain.RefrshTokenResponse{AcessToken: acessToken , Status: 200}


}

//  profile endpoint

func (u *UserUseCase) Profile(c context.Context, id string) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	user , err := u.UserRepository.FindUserByID(ctx , id)
	if err != nil {
		return domain.ErrorResponse{Message: "Error finding user" , Status: 500}
	}

	if id != user.ID.Hex() {
		return domain.ErrorResponse{Message: "Not Authorized" , Status: 400}
	}

	user.Password = ""
	user.RefreshToken = ""
	return domain.ProfileResponse{Message: "User found" , User: user ,  Status: 200}
}

func (u *UserUseCase) GetAllUsers(c context.Context , id string) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	// find user by id
	user , err := u.UserRepository.FindUserByID(ctx , id)

	if err != nil {
		return domain.ErrorResponse{Message: "Error finding user" , Status: 500}
	}

	if user.Role != "admin" {
		return domain.ErrorResponse{Message: "Not Authorized" , Status: 400}
	}

	allusers , err := u.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return domain.ErrorResponse{Message: "Error fetching users" , Status: 500}
	}

	return domain.AllUserResponse{Message: "Users found" , Status: 200 , Users: allusers} 

}
//  delete User

func (u *UserUseCase) DeleteUser(c context.Context, id string , adminid string) interface{} {
	ctx , cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// find user by id
	user , err := u.UserRepository.FindUserByID(ctx , adminid)

	if err != nil {
		return domain.ErrorResponse{Message: "Error finding user" , Status: 500}
	}

	if user.Role != "admin" {
		return domain.ErrorResponse{Message: "Not Authorized" , Status: 400}
	}

	// delete user

	err = u.UserRepository.DeleteUser(ctx , id)

	if err != nil {

		return domain.ErrorResponse{Message: "Error deleting user" , Status: 500}
	}

	return domain.SuccessResponse{Message: "User deleted successfully" , Status: 200}

}