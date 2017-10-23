package server

import (
	"config"
	"controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter for routing requests
func NewRouter() *echo.Echo {
	conf := config.GetConfig()
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	numberService := router.Group("/numberservice")
	version := numberService.Group("/" + conf.GetString("general.version"))
	number := version.Group("/number", middleware.BasicAuth(
		func(username, password string, c echo.Context) (bool, error) {
			if username == conf.GetString("auth.user") && password == conf.GetString("auth.password") {
				return true, nil
			}
			return false, nil
		}))

	number.GET("/:id/", controllers.Test)
	return router
}