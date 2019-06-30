package main

import (
	"github.com/kaz/flos-garden/collector"
	"github.com/kaz/flos-garden/database"
)

func main() {
	database.Init()
	collector.Init()
}
