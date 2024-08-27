package router

import (
	controller "loantracker/delivery/controller"
	"loantracker/infrastructure"
	"loantracker/mongo"
	"loantracker/repository"
	"loantracker/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func LoanRouter(config infrastructure.Config, DB mongo.Database, server *gin.RouterGroup) {

	loanrepo := repository.NewLoanRepository(DB, config.LoanCollection)
	userRepo := repository.NewUserRepository(DB, config.UserCollection)
	loanUsecase := usecase.NewLoanUsecase(loanrepo, time.Duration(config.ContextTimeout)*time.Second, &config , userRepo)

	loanController := controller.LoanController{ 
		LoanUseCase: loanUsecase,
	}

	authsecret := config.AccessTokenSecret
	server.Use(infrastructure.AuthMiddleware(authsecret))

	server.POST("loan" ,  loanController.CreateLoan)
	server.GET("loan/:id", loanController.GetSpecificLoan)
	server.GET("admin/loans" , loanController.GetAllLoans)
	server.PATCH("admin/loans/:id/status", loanController.UpdateLoanStatus)
	server.DELETE("admin/loans/:id", loanController.DeleteLoan)
	server.GET("admin/logs")


}