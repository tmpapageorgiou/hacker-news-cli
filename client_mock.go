package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type receiveEntry struct {
	Value []byte
	Error error
}

type MockJSONGet struct {
	Receive map[string]receiveEntry
	Count   int64
}

func NewMockJSONGet() *MockJSONGet {
	return &MockJSONGet{Receive: make(map[string]receiveEntry)}
}

func (m *MockJSONGet) Add(uri string, value string, err error) {
	m.Receive[uri] = receiveEntry{[]byte(value), err}
}

func (m *MockJSONGet) JSONGet(uri string, v interface{}) error {
	ret, ok := m.Receive[uri]
	atomic.AddInt64(&m.Count, 1)
	if !ok {
		return fmt.Errorf("uri not found in map! %s", uri)
	}

	err := json.Unmarshal(ret.Value, v)
	if err != nil {
		return err
	}

	return ret.Error
}
