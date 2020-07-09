package apigateway

import "github.com/labstack/echo/v4"

func Route() {
	e := echo.New()

	e.GET("/:service/manifest", getmainfest)
}
