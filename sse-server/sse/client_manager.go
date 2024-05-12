package sse

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

var mutex = &sync.Mutex{}

type Event struct {
	EventID string `json:"event_id"`
	Data    string `json:"data"`
	UserID  string `json:"user_id"`
	Time    string `json:"time"`
}

type ClientManager interface {
	AddClient(userID string) *Client
	RemoveClient(userID string)
	Broadcast(channel, message string)
}

type clientManager struct {
	clients map[string]*Client
}

func (c *clientManager) AddClient(userID string) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	existingClient, ok := c.clients[userID]
	if ok {
		atomic.AddInt64(&existingClient.TotalConnection, 1)
		return existingClient
	} else {
		client := &Client{
			ID:   userID,
			Send: make(chan string),
		}
		c.clients[userID] = client
		return client
	}
}

func (c *clientManager) RemoveClient(userID string) {
	mutex.Lock()
	defer mutex.Unlock()
	existingClient, ok := c.clients[userID]
	if ok {
		atomic.AddInt64(&existingClient.TotalConnection, -1)
		if existingClient.TotalConnection == 0 {
			delete(c.clients, userID)
		}
	}
}

func (c *clientManager) Broadcast(channel, message string) {
	var event Event
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		fmt.Println("error while unmarshal message", err)
		return
	}

	if client, ok := c.clients[event.UserID]; ok {
		// ignore compare channel
		// should check if there is any connection
		if client.TotalConnection > 0 {
			client.Send <- event.Data
		}
	}
}

type Client struct {
	ID              string
	Send            chan string
	TotalConnection int64
}

func NewClientManager() ClientManager {
	return &clientManager{
		clients: make(map[string]*Client),
	}
}
