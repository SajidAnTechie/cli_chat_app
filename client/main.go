package main

import (
	"bufio"
	"log"
	"os"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

func main() {

	opts := &socketio_client.Options{
		Transport: "websocket",
	}

	uri := "http://localhost:8000/socket.io/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func(s socketio_client.Client) {

		log.Printf("on connect\n")
	})

	client.On("message", func(msg string) {
		log.Printf("message:%v\n", msg)
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("chatMessage", command)
	}
}
