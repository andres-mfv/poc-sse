package main

import (
	"github.com/andres-mfv/sse-server/pubsub"
	"github.com/andres-mfv/sse-server/sse"
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

	clientManager := sse.NewClientManager()

	pubSubClient := pubsub.NewRedisClient(clientManager)

	// poc scope only have one channel
	go pubSubClient.Subscribe("sse_event")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/publish", func(c *gin.Context) {
		channel := c.PostForm("channel")
		message := c.PostForm("message")
		err := pubSubClient.Publish(channel, message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "published",
		})
	})

	r.GET("/sse", func(ctx *gin.Context) {
		// Set the headers required for Server-Sent Events
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")

		// get ctx user id param for poc only
		user_id := ctx.Param("user_id")
		user_id = "19"
		// add new client to client manager
		client := clientManager.AddClient(user_id)
		log.Println("create channel")
		// Continuously listen for messages to send to the client
		for {
			select {
			case msg := <-client.Send:
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
