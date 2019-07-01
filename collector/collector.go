package collector

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/kaz/flos-garden/collector/lifeline"

	"github.com/kaz/flos-garden/collector/bookshelf"
	"github.com/kaz/flos-garden/common"
	"github.com/kaz/flos-garden/database"
)

var (
	logger = log.New(os.Stdout, "[collector] ", log.Ltime)

	mu      = sync.RWMutex{}
	workers = map[string]context.CancelFunc{}
)

func Init() {
	bookshelf.Init()

	if _, err := database.DB().Exec("CREATE TABLE IF NOT EXISTS instances (host TEXT, bastion BOOL, PRIMARY KEY(host(128)))" + common.TABLE_OPTION); err != nil {
		panic(err)
	}

	rows, err := database.DB().Query("SELECT * FROM instances")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var host string
		var bastion bool
		if err := rows.Scan(&host, &bastion); err != nil {
			panic(err)
		}

		if bastion {
			common.RegisterBastion(host)
		} else {
			go runWorker(host)
		}
	}
}

func runWorker(addr string) {
	ctx, cancel := context.WithCancel(context.Background())

	mu.Lock()
	workers[addr] = cancel
	mu.Unlock()

	if err := lifeline.RunLifelineCollector(ctx, addr); err != nil {
		logger.Printf("failed to start lifeline collector: %v (host=%s)\n", err, addr)
		return
	}
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
