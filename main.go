package main

import (
	"cli_chat_app/cmd"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {

	cmd.Execute()

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {

		fmt.Println("connected successfull")

		s.SetContext("")
		s.Join("chat")
		s.Emit("message", "Welcom to the chat")

		server.BroadcastToRoom("", "chat", "message", "user with id"+s.ID()+"join the chat")

		fmt.Println("connected:")

		return nil
	})

	server.OnEvent("/", "chat message", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		server.BroadcastToRoom("", "chat", "chat message", msg)
		fmt.Println("notice:", msg)

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
