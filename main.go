package main

import (
	"net/http"

	"github.com/kaz/flos-garden/collector"
	"github.com/kaz/flos-garden/database"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(http.StatusInternalServerError, map[string]error{"error": err})
	}

	database.Init()
	collector.Init()

	api := e.Group("/api")
	collector.RegisterHandler(api.Group("/collector"))

	e.Logger.Fatal(e.Start(":9000"))
}
