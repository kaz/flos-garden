package database

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	QueryRequest struct {
		Query string
	}
)

func RegisterHandler(g *echo.Group) {
	g.POST("/query", postQuery)
}

func postQuery(c echo.Context) error {
	var req QueryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	rows, err := Query(req.Query)
	if err != nil {
		return err
	}

	resp := []map[string]interface{}{}
	for rows.Next() {
		row := map[string]interface{}{}
		if err := rows.MapScan(row); err != nil {
			return err
		}

		for k, v := range row {
			byteValue, ok := v.([]byte)
			if ok {
				row[k] = string(byteValue)
			}
		}

		resp = append(resp, row)
	}

	return c.JSON(http.StatusOK, resp)
}
