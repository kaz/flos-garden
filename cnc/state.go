package cnc

import (
	"fmt"
	"net/http"

	"github.com/kaz/flos-hortus/common"
	"github.com/kaz/flos/state"
	"github.com/labstack/echo/v4"
)

func getState(c echo.Context) error {
	resp, err := common.Request(http.MethodGet, c.Param("host"), "/state", nil, nil)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to get state: %v\n", msg)
	}

	var state state.State
	if err := common.ReadBody(resp, &state); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}

	return c.JSON(http.StatusOK, state)
}

func putState(c echo.Context) error {
	var state state.State
	if err := c.Bind(&state); err != nil {
		return fmt.Errorf("failed to bind request body: %v\n", err)
	}

	resp, err := common.Request(http.MethodPut, c.Param("host"), "/state", state, nil)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to get state: %v\n", msg)
	}

	return c.NoContent(http.StatusOK)
}
