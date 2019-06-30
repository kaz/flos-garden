package collector

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/kaz/flos-garden/database"
)

const (
	INSTANCE_CHECK_SEC = 10

	TABLE_OPTION = "CHARACTER SET utf8mb4"
)

var (
	logger = log.New(os.Stdout, "[collector] ", log.Ltime)

	mu      = sync.RWMutex{}
	workers = map[string]context.CancelFunc{}
)

func Init() {
	if _, err := database.Exec("CREATE TABLE IF NOT EXISTS instances (id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, addr TEXT)" + TABLE_OPTION); err != nil {
		panic(err)
	}

	go runMaster()
}

func runMaster() {
	for {
		rows, err := database.Query("SELECT addr FROM instances")
		if err != nil {
			logger.Printf("query failed: %v\n", err)
			continue
		}

		addrs := map[string]interface{}{}
		for rows.Next() {
			var addr string
			if err := rows.Scan(&addr); err != nil {
				logger.Printf("scan failed: %v\n", err)
				continue
			}

			addrs[addr] = nil
			go runWorker(addr)
		}

		mu.RLock()
		for addr, cancel := range workers {
			if _, ok := addrs[addr]; !ok {
				cancel()
				delete(workers, addr)
			}
		}
		mu.RUnlock()

		time.Sleep(INSTANCE_CHECK_SEC * time.Second)
	}
}

func runWorker(addr string) {
	_, cancel := context.WithCancel(context.Background())

	mu.Lock()
	workers[addr] = cancel
	mu.Unlock()

}
