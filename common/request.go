package common

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kaz/flos/messaging"
)

var (
	bastion string
)

func RegisterBastion(host string) {
	bastion = host
}

func Request(method string, host string, path string, data interface{}, ctx context.Context) (*http.Response, error) {
	endpoint := "http://" + host + REMOTE_LISTEN + path
	authorization := ""

	if bastion != "" {
		payload, err := messaging.Encode(host + REMOTE_LISTEN)
		if err != nil {
			return nil, fmt.Errorf("failed to create payload: %v\n", err)
		}

		endpoint = "http://" + bastion + REMOTE_LISTEN + path
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

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v\n", err)
	}

	return resp, nil
}

func ReadBody(resp *http.Response, ptr interface{}) error {
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
