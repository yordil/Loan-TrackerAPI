package router

import (
	"loantracker/infrastructure"
	"loantracker/mongo"

	"github.com/gin-gonic/gin"
)

func Router(server *gin.RouterGroup, config *infrastructure.Config, DB mongo.Database) {


	// setting up new router group for user

	userRouter := server.Group("")

	// call the function from user_router.go
	UserRouter(*config, DB , userRouter)

	// adminRouter := server.Group("")

	// AdminRouter(*config, DB, adminRouter)

	// adminRouter := server.Group("/admin")

	// AdminRouter(*config, DB, adminRouter)
}