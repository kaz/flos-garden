package bookshelf

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kaz/flos-garden/database"
	"github.com/kaz/flos/libra/bookshelf"
	"github.com/kaz/flos/messaging"
)

const (
	TABLE_OPTION  = "CHARACTER SET utf8mb4"
	REMOTE_LISTEN = ":10239"
	COLLECT_SEC   = 8
)

type (
	collector struct {
		logger *log.Logger

		cur   uint64
		table string

		name string
		host string
		path string

		ctx context.Context
	}
)

var (
	tz      *time.Location
	bastion string
)

func Init() {
	var err error
	if tz, err = time.LoadLocation("Asia/Tokyo"); err != nil {
		panic(err)
	}

	if _, err = database.DB().Exec("CREATE TABLE IF NOT EXISTS bookshelf_cursor (name TEXT, host TEXT, cur BIGINT UNSIGNED, PRIMARY KEY(name(128), host(128)))" + TABLE_OPTION); err != nil {
		panic(err)
	}
}

func RegisterBastion(host string) {
	bastion = host
}

func newBookshelfCollector(ctx context.Context, name string, host string, path string, contentType string) (*collector, error) {
	logger := log.New(os.Stdout, "[bookshelf name="+name+" host="+host+"] ", log.Ltime)

	table := "bookshelf_data_" + name
	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS " + table + " (host TEXT, remote_id BIGINT UNSIGNED, series TEXT, contents " + contentType + ", created DATETIME(6), PRIMARY KEY(host(128), remote_id), KEY(series(128), created))" + TABLE_OPTION); err != nil {
		return nil, err
	}

	var cur uint64
	if err := database.DB().QueryRow("SELECT cur FROM bookshelf_cursor WHERE name = ? AND host = ?", name, host).Scan(&cur); err != nil {
		logger.Println("No cursor found. Set to 0.")
		cur = 0
	}

	return &collector{logger, cur, table, name, host, path, ctx}, nil
}

func (c *collector) Collect() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if err := c.collect(); err != nil {
				c.logger.Println("collect failed:", err)
			} else {
				c.logger.Println("collected")
			}
		}
		time.Sleep(COLLECT_SEC * time.Second)
	}
}

func (c *collector) doRequest(method string, data uint64) (*http.Response, error) {
	endpoint := "http://" + c.host + REMOTE_LISTEN + c.path
	authorization := ""

	if bastion != "" {
		payload, err := messaging.Encode(c.host + REMOTE_LISTEN)
		if err != nil {
			return nil, fmt.Errorf("failed to create payload: %v\n", err)
		}

		endpoint = "http://" + bastion + REMOTE_LISTEN + c.path
		authorization = "bearer " + base64.StdEncoding.EncodeToString(payload)
	}

	payload, err := messaging.Encode(float64(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create payload: %v\n", err)
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v\n", err)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(c.ctx))
	if err != nil {
		return nil, fmt.Errorf("request failed: %v\n", err)
	}

	return resp, nil
}

func (c *collector) collect() error {
	resp, err := c.doRequest(http.MethodPatch, c.cur)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		var msg string
		if err := messaging.Decode(respBody, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}

		return fmt.Errorf("failed to get books: %v\n", msg)
	}

	var books []*bookshelf.Book
	if err := messaging.Decode(respBody, &books); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}

	stmt, err := database.DB().PrepareNamedContext(c.ctx, "REPLACE INTO "+c.table+" VALUES (:host, :remote_id, :series, :contents, :created)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v\n", err)
	}

	var curMax uint64
	for _, book := range books {
		_, err := stmt.ExecContext(c.ctx, map[string]interface{}{
			"host":      c.host,
			"remote_id": book.ID,
			"series":    book.Series,
			"contents":  book.Contents,
			"created":   time.Unix(0, book.Timestamp).In(tz),
		})
		if err != nil {
			return fmt.Errorf("failed to insert record: %v\n", err)
		}
		if book.ID > curMax {
			curMax = book.ID
		}
	}

	delResp, err := c.doRequest(http.MethodDelete, curMax)
	if err != nil {
		return fmt.Errorf("request failed: %v\n", err)
	}
	defer delResp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v\n", err)
		}

		var msg string
		if err := messaging.Decode(respBody, &msg); err != nil {
			return fmt.Errorf("failed to decode resp body: %v\n", err)
		}

		return fmt.Errorf("failed to delete books: %v\n", string(respBody))
	}

	if _, err := database.DB().ExecContext(c.ctx, "REPLACE INTO bookshelf_cursor VALUES (?, ?, ?)", c.name, c.host, curMax+1); err != nil {
		return fmt.Errorf("failed to update cursor: %v\n", err)
	}

	c.cur = curMax + 1
	return nil
}
