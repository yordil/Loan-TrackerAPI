package controller

import (
	"loantracker/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)


func HandleResponse(c *gin.Context, response interface{}) {

	switch res := response.(type) {
	case domain.SuccessResponse:
		c.JSON(http.StatusOK, res)
	case domain.ErrorResponse:
		c.JSON(res.Status, res)
	case domain.LoginResponse:
		c.JSON(http.StatusOK, res)
	case domain.AllUserResponse:
		c.JSON(http.StatusOK, res)
	case domain.ProfileResponse:
		c.JSON(http.StatusOK, res)
	case domain.RefrshTokenResponse:
		c.JSON(http.StatusOK, res)
	
	default:
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Internal Server Error", Status: 500})
	}
}
