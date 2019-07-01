package bookshelf

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/kaz/flos-garden/database"
)

const (
	TABLE_OPTION = "CHARACTER SET utf8mb4"

	COLLECT_SEC = 8
)

type (
	collector struct {
		logger *log.Logger

		cur      int
		table    string
		endpoint string
	}
)

func Init() {
	if _, err := database.Exec("CREATE TABLE IF NOT EXISTS bookshelf_cursor (name TEXT, host TEXT, cur INT UNSIGNED, PRIMARY KEY(name(128), host(128)))" + TABLE_OPTION); err != nil {
		panic(err)
	}
}

func newBookshelfCollector(name string, host string, path string) (*collector, error) {
	logger := log.New(os.Stdout, "[bookshelf name="+name+" host="+host+"] ", log.Ltime)

	table := "bookshelf_data_" + name
	if _, err := database.Exec("CREATE TABLE IF NOT EXISTS " + table + " (host TEXT, remote_id INT UNSIGNED, series TEXT, contents TEXT, created DATETIME, PRIMARY KEY(host(128), remote_id), KEY(series(128), created))" + TABLE_OPTION); err != nil {
		return nil, err
	}

	var cur int
	if err := database.QueryRow("SELECT cur FROM bookshelf_cursor WHERE name = ? AND host = ?", name, host).Scan(&cur); err != nil {
		logger.Println("No cursor found. Set to 0.")
		cur = 0
	}

	return &collector{logger, cur, table, "http://" + host + ":10239" + path}, nil
}

func (c *collector) Collect(ctx context.Context) {
	for {
		c.logger.Println("Collecting ...")

		time.Sleep(COLLECT_SEC * time.Second)
	}
}
