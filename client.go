package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JSONGetter interface {
	JSONGet(string, interface{}) error
}

type HTTPJSONClient struct {
	client *http.Client
}

func NewHTTPJSONClient() *HTTPJSONClient {
	return &HTTPJSONClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h *HTTPJSONClient) JSONGet(uri string, v interface{}) error {
	resp, err := h.client.Get(uri)
	if err != nil {
		return fmt.Errorf("failed to request %s: %v", uri, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get from %s: status code [%d]", uri, resp.StatusCode)
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body response: %v", err)
	}

	err = json.Unmarshal(responseBytes, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return err
}
