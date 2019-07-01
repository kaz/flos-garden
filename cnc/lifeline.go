package cnc

import (
	"fmt"
	"net/http"

	"github.com/kaz/flos-garden/common"
	"github.com/labstack/echo/v4"
)

func postShell(c echo.Context) error {
	var shell string
	if err := c.Bind(&shell); err != nil {
		return fmt.Errorf("failed to bind request body: %v\n", err)
	}

	resp, err := common.Request(http.MethodPost, c.Param("host"), "/lifeline/shell", shell, nil)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to exec command: %v\n", msg)
	}

	var out string
	if err := common.ReadBody(resp, &out); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}

	return c.JSON(http.StatusOK, out)
}
