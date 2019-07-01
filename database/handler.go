package database

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandler(g *echo.Group) {
	g.POST("/query", postQuery)
}

func postQuery(c echo.Context) error {
	sql, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	rows, err := DB().Queryx(string(sql))
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
