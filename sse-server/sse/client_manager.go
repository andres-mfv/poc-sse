package sse

import (
	"encoding/json"
	"fmt"
	"sync"
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
	RemoveClient(client *Client)
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

func (c *clientManager) RemoveClient(client *Client) {
	//TODO implement me
	panic("implement me")
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
		client.Send <- event.Data
	}
}

type Client struct {
	ID   string
	Send chan string
}

func NewClientManager() ClientManager {
	return &clientManager{
		clients: make(map[string]*Client),
	}
}
