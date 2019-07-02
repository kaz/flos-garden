package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kaz/flos-hortus/cnc"
	"github.com/kaz/flos-hortus/collector"
	"github.com/kaz/flos-hortus/database"
	"github.com/kaz/flos/messaging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	BASIC_USER = "k4mu1"
	BASIC_PASS = "k4mu1"
)

var (
	logger = log.New(os.Stdout, "[http] ", log.Ltime)
)

func main() {
	messaging.Init()
	collector.Init()

	e := echo.New()
	e.HideBanner = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		logger.Println(err)

		status := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			status = he.Code
		}

		c.JSON(status, map[string]string{"error": err.Error()})
	}

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == BASIC_USER && password == BASIC_PASS, nil
	}))

	api := e.Group("/api")
	cnc.RegisterHandler(api.Group("/cnc"))
	database.RegisterHandler(api.Group("/database"))
	collector.RegisterHandler(api.Group("/collector"))

	e.Static("/", "./view/dist")
	e.Logger.Fatal(e.Start(":9000"))
}
