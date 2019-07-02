package cnc

import (
	"fmt"
	"net/http"

	"github.com/kaz/flos-hortus/common"
	"github.com/labstack/echo/v4"
)

func postPower(c echo.Context) error {
	var power string
	if err := c.Bind(&power); err != nil {
		return fmt.Errorf("failed to bind request body: %v\n", err)
	}

	resp, err := common.Request(http.MethodPost, c.Param("host"), "/power", power, nil)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to switch power: %v\n", msg)
	}

	return c.NoContent(http.StatusOK)
}
