package common

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kaz/flos/messaging"
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

var (
	bastion string
)

func RegisterBastion(host string) {
	bastion = host
}

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
	endpoint := "http://" + c.Host + REMOTE_LISTEN + c.Path
	authorization := ""

	if bastion != "" {
		payload, err := messaging.Encode(c.Host + REMOTE_LISTEN)
		if err != nil {
			return nil, fmt.Errorf("failed to create payload: %v\n", err)
		}

		endpoint = "http://" + bastion + REMOTE_LISTEN + c.Path
		authorization = "bearer " + base64.StdEncoding.EncodeToString(payload)
	}

	payload, err := messaging.Encode(data)
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

	resp, err := http.DefaultClient.Do(req.WithContext(c.Context))
	if err != nil {
		return nil, fmt.Errorf("request failed: %v\n", err)
	}

	return resp, nil
}

func (c *Collector) ReadBody(resp *http.Response, ptr interface{}) error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v\n", err)
	}
	resp.Body.Close()

	if err := messaging.Decode(respBody, ptr); err != nil {
		return fmt.Errorf("failed to decode resp body: %v\n", err)
	}
	return nil
}
