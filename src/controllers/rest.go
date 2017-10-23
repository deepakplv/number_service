package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

// Sample controller
func Test(c echo.Context) error {
	id := c.Param("id")

	var response echo.Map
	response = echo.Map{
		"data":   "ID is - "+id,
		"status": "success",
	}

	return c.JSON(http.StatusOK, response)

}