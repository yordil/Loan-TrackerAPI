package main

import (
	"fmt"
	"loantracker/delivery/router"
	"loantracker/infrastructure"

	xss "github.com/araujo88/gin-gonic-xss-middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	var xssMdlwr xss.XssMw
	server.Use(xssMdlwr.RemoveXss())

	


	config, err := infrastructure.LoadEnv()
	if err != nil {
		fmt.Print("Error in env.load")
	}
	// fmt.Print(config)
	DB, client, err := infrastructure.ConnectDB(config.DatabaseUrl, config.DbName)

	if err != nil {
		fmt.Print("Error in connectDB")
	}
	defer infrastructure.CloseDB(client)
	router.Router(server.Group("api/v0"), config, DB)
	server.Run(fmt.Sprintf(":%d", config.Port))

}
