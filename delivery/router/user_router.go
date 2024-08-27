package router

import (
	controller "loantracker/delivery/controller"
	"loantracker/mongo"
	"loantracker/repository"
	"loantracker/usecase"
	"time"

	"github.com/gin-gonic/gin"

	"loantracker/infrastructure"
)



func UserRouter(config infrastructure.Config , DB mongo.Database ,  server *gin.RouterGroup) {

	

	passwordService := infrastructure.NewPasswordService()

	userRepo := repository.NewUserRepository(DB, config.UserCollection)
	userUse := usecase.NewUserUseCase(userRepo ,  time.Duration(config.ContextTimeout)*time.Second , passwordService , config)

	userController := controller.UserController{
		UserUseCase: userUse,
	}

	authsecret := config.AccessTokenSecret
	server.POST("user/create", userController.RegisterUser)
	server.GET("user/verify-email" , userController.VerifyUser)
	server.POST("user/login" , userController.Login)
	server.POST("user/password-reset" , userController.ForgotPassword)
	server.POST("user/token/refresh" , infrastructure.AuthMiddleware(authsecret),  userController.RefreshToken)
	server.GET("user/profile" , infrastructure.AuthMiddleware(authsecret) ,  userController.Profile)
	

	server.GET("admin/user" , infrastructure.AuthMiddleware(authsecret) ,  userController.GetAllUsers)
	server.DELETE("admin/user/:id" , infrastructure.AuthMiddleware(authsecret) , userController.DeleteUser)

}