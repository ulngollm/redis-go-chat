package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func readRequest(inputChannel string) {
	for {
		var inputMsg string
		reader := bufio.NewReader(os.Stdin)
		inputMsg, err := reader.ReadString('\n')

		if inputMsg[0] == '>' {
			continue
		}

		if err != nil {
			fmt.Println(err)
		}
		rdb.Publish(ctx, inputChannel, inputMsg)
	}
}

func printResponse(subscriber *redis.PubSub) {
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		if msg.Payload[0] != '>' {
			continue
		}
		fmt.Println(msg.Payload)
	}
}

func main() {
	inputChannel := faker.Word()
	rdb.Publish(ctx, "openedChat", inputChannel)
	fmt.Println(inputChannel)
	subscriber := rdb.Subscribe(ctx, inputChannel)
	go printResponse(subscriber)
	readRequest(inputChannel)
}
