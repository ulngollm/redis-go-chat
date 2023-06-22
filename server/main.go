package main

import (
	"context"
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

const inputChannel = "openedChat"

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func handleConversation(channel string) {
	subscriber := rdb.Subscribe(ctx, channel)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		if msg.Payload[0] == '>' {
			continue
		}

		responseMsg := fmt.Sprintf(">%s", faker.Paragraph())
		err = rdb.Publish(ctx, channel, responseMsg).Err()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	subscriber := rdb.Subscribe(ctx, inputChannel)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		channel := msg.Payload
		fmt.Println(msg.Payload)

		go handleConversation(channel)
	}
}
