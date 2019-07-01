package action

import "github.com/labstack/echo/v4"

func RegisterHandler(g *echo.Group) {
	g.GET("/state/:host", getState)
	g.PUT("/state/:host", putState)
}
