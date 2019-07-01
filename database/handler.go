package database

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterHandler(g *echo.Group) {
	g.GET("/blob/:host/:id", getBLOB)
	g.POST("/query", postQuery)
}

func getBLOB(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse request body: %v\n", err)
	}

	var path string
	var data []byte
	if err := DB().QueryRow("SELECT series, contents FROM bookshelf_data_archive WHERE host = ? AND remote_id = ?", c.Param("host"), id).Scan(&path, &data); err != nil {
		return fmt.Errorf("failed to get BLOB: %v\n", err)
	}

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(path)))
	return c.Blob(http.StatusOK, "application/octet-stream", data)
}

func postQuery(c echo.Context) error {
	var sql string
	if err := c.Bind(&sql); err != nil {
		return fmt.Errorf("failed to bind request body: %v\n", err)
	}

	rows, err := DB().Queryx(sql)
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
