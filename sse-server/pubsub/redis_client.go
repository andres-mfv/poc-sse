package pubsub

import (
	"context"
	"fmt"
	"github.com/andres-mfv/sse-server/sse"
	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:16379", // Redis server address
	Username: "default",
	Password: "moneyforward123",
})

func init() {
	rs, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(rs)
}

type PubSub interface {
	Publish(channel string, message interface{}) error
	Subscribe(channel string)
}

type redisClient struct {
	pubsub        *redis.Client
	clientManager sse.ClientManager
}

func (r *redisClient) Publish(channel string, message interface{}) error {
	return r.pubsub.Publish(context.TODO(), channel, message).Err()
}

func (r *redisClient) Subscribe(channel string) {
	subscribe := r.pubsub.Subscribe(context.TODO(), channel)
	ch := subscribe.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		r.clientManager.Broadcast(msg.Channel, msg.Payload)
	}
}

func NewRedisClient(clientManager sse.ClientManager) PubSub {
	return &redisClient{
		pubsub:        rdb,
		clientManager: clientManager,
	}
}
