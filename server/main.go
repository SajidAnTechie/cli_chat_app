package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {

	server := socketio.NewServer(nil)

	server.OnConnect("connection", func(s socketio.Conn) error {

		fmt.Println("connected successfull")

		s.Emit("message", "Welcom to the chat")

		fmt.Println("===Start Charting====")

		return nil
	})
	server.BroadcastToRoom("", "chat", "message", "user join the chat")

	server.OnEvent("/", "chatMessage", func(s socketio.Conn, msg string) {

		//s.BroadcastToRoom("", "chat", "message", msg)

		s.Emit("message", msg)

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

		fmt.Println("User with id left the chat", s.ID())
		fmt.Println("closed", reason)

	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
