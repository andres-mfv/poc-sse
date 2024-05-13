package sse

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	RemoveClient(clientID string)
	Broadcast(channel, message string)
}

type clientManager struct {
	clients map[string]map[string]*Client
}

func (c *clientManager) AddClient(userID string) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	clientID := uuid.New().String()
	fmt.Println("Add new client ", clientID, " for user ", userID)
	client := &Client{
		ID:   clientID,
		Send: make(chan string),
	}
	mapOfClients, ok := c.clients[userID]
	if !ok {
		mmClient := map[string]*Client{
			clientID: client,
		}
		c.clients[userID] = mmClient
	} else {
		mapOfClients[clientID] = client
	}
	return client
}

func (c *clientManager) RemoveClient(clientID string) {
	fmt.Println("Remove client", clientID)
	mutex.Lock()
	defer mutex.Unlock()
	for userID, mapOfClients := range c.clients {
		if _, ok := mapOfClients[clientID]; ok {
			delete(mapOfClients, clientID)
			if len(mapOfClients) == 0 {
				delete(c.clients, userID)
			}
			break
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

	for _, clients := range c.clients[event.UserID] {
		clients.Send <- message
	}
}

type Client struct {
	ID   string
	Send chan string
}

func NewClientManager() ClientManager {
	return &clientManager{
		clients: make(map[string]map[string]*Client),
	}
}
