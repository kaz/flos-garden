package bookshelf

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kaz/flos-garden/common"
	"github.com/kaz/flos-garden/database"
	"github.com/kaz/flos/libra/bookshelf"
)

type (
	collector struct {
		*common.Collector

		name  string
		table string
		cur   uint64
	}
)

func Init() {
	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS bookshelf_cursor (host TEXT, name TEXT, cur BIGINT UNSIGNED, PRIMARY KEY(host(128), name(128)))" + common.TABLE_OPTION); err != nil {
		panic(err)
	}
}

func newBookshelfCollector(ctx context.Context, name string, host string, path string, contentType string) (*collector, error) {
	logger := log.New(os.Stdout, "[bookshelf name="+name+" host="+host+"] ", log.Ltime)

	table := "bookshelf_data_" + name
	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS " + table + " (host TEXT, remote_id BIGINT UNSIGNED, series TEXT, contents " + contentType + ", created DATETIME(6), PRIMARY KEY(host(128), remote_id), KEY(series(128), created))" + common.TABLE_OPTION); err != nil {
		return nil, err
	}

	var cur uint64
	if err := database.DB().QueryRow("SELECT cur FROM bookshelf_cursor WHERE host = ? AND name = ?", host, name).Scan(&cur); err != nil {
		logger.Println("No cursor found. Set to 0.")
		cur = 0
	}

	c := &collector{
		Collector: &common.Collector{
			Context: ctx,
			Logger:  logger,
			Host:    host,
			Path:    path,
		},

		name:  name,
		table: table,
		cur:   cur,
	}

	c.RegisterCollectFunc(c.collect)
	return c, nil
}

func (c *collector) collect() error {
	resp, err := c.DoRequest(http.MethodPatch, float64(c.cur))
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to get books: %v\n", msg)
	}

	var books []*bookshelf.Book
	if err := common.ReadBody(resp, &books); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}

	stmt, err := database.DB().PrepareNamedContext(c.Context, "REPLACE INTO "+c.table+" VALUES (:host, :remote_id, :series, :contents, :created)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v\n", err)
	}

	var curMax uint64
	for _, book := range books {
		_, err := stmt.ExecContext(c.Context, map[string]interface{}{
			"host":      c.Host,
			"remote_id": book.ID,
			"series":    book.Series,
			"contents":  book.Contents,
			"created":   common.Time(book.Timestamp),
		})
		if err != nil {
			return fmt.Errorf("failed to insert record: %v\n", err)
		}
		if book.ID > curMax {
			curMax = book.ID
		}
	}

	delResp, err := c.DoRequest(http.MethodDelete, float64(curMax))
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}
	defer delResp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := common.ReadBody(resp, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}
		return fmt.Errorf("failed to get books: %v\n", msg)
	}

	if _, err := database.DB().ExecContext(c.Context, "REPLACE INTO bookshelf_cursor VALUES (?, ?, ?)", c.Host, c.name, curMax+1); err != nil {
		return fmt.Errorf("failed to update cursor: %v\n", err)
	}

	c.cur = curMax + 1
	return nil
}
