package collector

import (
	"context"
	"log"
	"os"
	"sync"

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
	if _, err := database.Exec("CREATE TABLE IF NOT EXISTS instances (id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, addr TEXT UNIQUE, bastion BOOL)" + TABLE_OPTION); err != nil {
		panic(err)
	}

	rows, err := database.Query("SELECT addr FROM instances")
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
	_, cancel := context.WithCancel(context.Background())

	mu.Lock()
	workers[addr] = cancel
	mu.Unlock()

	logger.Println("worker started:", addr)
}
