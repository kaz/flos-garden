package collector

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/kaz/flos-garden/collector/bookshelf"
	"github.com/kaz/flos-garden/database"
)

const (
	TABLE_OPTION = "CHARACTER SET utf8mb4"
)

var (
	logger = log.New(os.Stdout, "[collector] ", log.Ltime)

	mu      = sync.RWMutex{}
	workers = map[string]context.CancelFunc{}
)

func Init() {
	bookshelf.Init()

	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS instances (host TEXT, bastion BOOL, PRIMARY KEY(host(128)))" + TABLE_OPTION); err != nil {
		panic(err)
	}

	rows, err := database.DB().Query("SELECT host FROM instances")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var addr string
		if err := rows.Scan(&addr); err != nil {
			panic(err)
		}

		go runWorker(addr)
	}
}

func runWorker(addr string) {
	ctx, cancel := context.WithCancel(context.Background())

	mu.Lock()
	workers[addr] = cancel
	mu.Unlock()

	if err := bookshelf.RunLibraCollector(ctx, addr); err != nil {
		logger.Printf("failed to start libra collector: %v (host=%s)\n", err, addr)
		return
	}
	if err := bookshelf.RunArchiveCollector(ctx, addr); err != nil {
		logger.Printf("failed to start archive collector: %v (host=%s)\n", err, addr)
		return
	}
	logger.Printf("worker started (host=%s)\n", addr)
}
