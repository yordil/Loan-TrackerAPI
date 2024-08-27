package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempity" json:"id" `
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Role    string             `json:"role"`	
	IsVerified bool 			`json:"is_verified"`
	ResetPasswordToken   string               `json:"reset_password_token"`
	ResetPasswordExpires time.Time            `json:"reset_password_expires"`
	VerificationToken  string 					`json:"verification"`
	VerificationExpiration  time.Time 					`json:"expiration"`
	RefreshToken string 			`json:"refresh_token"`
}



type UserRepository interface { 
	CreateUser(c context.Context , user User) (User, error)
	FindUserByEmail(c context.Context ,   email string) (User, error)
	UpdateUser(c context.Context ,user  User) (User, error)
	GetAllUsers(c context.Context) ([]User, error)
	// FindUserByUsername(username string) (User, error)
	// FindUserByID(id string) (User, error)
	// PasswordReset(user User) (User, error)
	// GetAllUsers() ([]User, error)
	// DeleteUser(id string) error
	// VerifyOTP(otp string) (bool, error)
	DeleteUser(c context.Context , id string) (error)
	FindUserByID(c context.Context , id string) (User, error)
	FindUserByResetToken(c context.Context , token string) (User, error)
	SetResetToken(c context.Context, email ForgotPasswordRequest, token string, expiration time.Time) (User, error)
	FindUserByVerificationToken(c context.Context , token string) (User, error)

}

type UserUseCase interface { 
	CreateUser(c context.Context ,user User ) interface{}
	VeryfyUser(c context.Context , token string)interface {}
	ForgotPassword(c context.Context, email ForgotPasswordRequest) interface{}
	ResetPassword(c context.Context, password ResetPasswordRequest, token string) interface{}
	RefreshToken(c context.Context, id string , refreshToken RefreshTokenRequest) interface{}
	Profile(c context.Context, id string) interface{}
	GetAllUsers(c context.Context , id string) interface{}
	DeleteUser(c context.Context, id string , adminid string) interface{}
	Login(c context.Context, login LoginRequest) interface{}
	// FindUserByEmail(email string) (User, error)
	// FindUserByUsername(username string) (User, error)
	// FindUserByID(id string) (User, error)
	// PasswordReset(user User) (User, error)
	// GetAllUsers() ([]User, error)
	// DeleteUser(id string) error
	// VerifyOTP(otp string) (bool, error)


}