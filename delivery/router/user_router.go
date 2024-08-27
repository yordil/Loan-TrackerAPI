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
	userUse := usecase.NewUserUseCase(userRepo ,  time.Duration(config.ContextTimeout)*time.Second , passwordService)

	userController := controller.UserController{
		UserUseCase: userUse,
	}


	server.POST("user/create", userController.RegisterUser)
	server.GET("user/verify-email" , userController.VerifyUser)
	server.POST("user/login" , userController.Login)
	server.POST("user/password-reset" , userController.ForgotPassword)
	server.POST("user/token/refresh" , infrastructure.AuthMiddleware(),  userController.RefreshToken)
	server.POST("user/profile" , infrastructure.AuthMiddleware() ,  userController.Profile)
	

	server.GET("admin/user" , infrastructure.AuthMiddleware() ,  userController.GetAllUsers)
	server.DELETE("admin/user/:id" , infrastructure.AuthMiddleware() , userController.DeleteUser)




	
}