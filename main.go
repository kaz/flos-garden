package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kaz/flos-garden/cnc"
	"github.com/kaz/flos-garden/collector"
	"github.com/kaz/flos-garden/database"
	"github.com/kaz/flos/messaging"
	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[http] ", log.Ltime)
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	messaging.Init()
	collector.Init()

	api := e.Group("/api")
	cnc.RegisterHandler(api.Group("/cnc"))
	database.RegisterHandler(api.Group("/database"))
	collector.RegisterHandler(api.Group("/collector"))

	e.Static("/", "./view/dist")
	e.Logger.Fatal(e.Start(":9000"))
}
