package common

import (
	"context"
	"log"
	"net/http"
	"time"
)

type (
	Collector struct {
		Context context.Context
		Logger  *log.Logger

		Host string
		Path string

		collect func() error
	}
)

func (c *Collector) RegisterCollectFunc(collectFunc func() error) {
	c.collect = collectFunc
}

func (c *Collector) Collect() {
	for {
		select {
		case <-c.Context.Done():
			return
		default:
			if err := c.collect(); err != nil {
				c.Logger.Println("collect failed:", err)
			} else {
				c.Logger.Println("collected")
			}
		}
		time.Sleep(COLLECT_CYCLE_SEC * time.Second)
	}
}

func (c *Collector) DoRequest(method string, data interface{}) (*http.Response, error) {
	return Request(method, c.Host, c.Path, data, c.Context)
}
