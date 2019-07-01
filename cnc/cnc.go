package cnc

import "github.com/labstack/echo/v4"

func RegisterHandler(g *echo.Group) {
	g.GET("/:host/state", getState)
	g.PUT("/:host/state", putState)
}
