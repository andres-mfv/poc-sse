package main

import (
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/sse", func(ctx *gin.Context) {
		// Set the headers required for Server-Sent Events
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")

		// Create a new channel for this client
		clientChan := make(chan string)
		log.Println("create channel")
		//defer close(clientChan)
		go func() {
			// every 5 seconds send a message to clientChan
			for {
				clientChan <- "Hello, world!"
				time.Sleep(1 * time.Second)
			}
		}()
		// Continuously listen for messages to send to the client
		for {
			select {
			case msg := <-clientChan:
				// Write the message to the response writer
				ctx.SSEvent("message", msg)
				ctx.Writer.Flush()
				log.Println("send message to client")
			case <-ctx.Writer.CloseNotify():
				// If the connection is closed, return
				return
			}
		}
	})

	err := r.Run()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
