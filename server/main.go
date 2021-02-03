package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type user struct {
	name     string
	roomName string
}

func main() {

	mapMake := make(map[string]*user)

	userJoin := func(userID string, joinedRoom string, userName string) {

		mapMake[userID] = &user{name: userName, roomName: joinedRoom}

	}
	getJoinedUserDetails := func(userID string) *user {
		return mapMake[userID]
	}

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {

		fmt.Println("connected successfull")

		s.Emit("message", "Welcom to the chat")

		return nil
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomName string, userName string) {

		userJoin(s.ID(), roomName, userName)

		s.Join(roomName)

		server.BroadcastToRoom("/", roomName, "message", userName+" join the chat")

	})

	server.OnEvent("/", "chatMessage", func(s socketio.Conn, msg string) {

		getJoinedUserDetails := getJoinedUserDetails(s.ID())

		server.BroadcastToRoom("/", getJoinedUserDetails.roomName, "message", getJoinedUserDetails.name+":  "+msg)

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

		s.Emit("leaveRoom", mapMake[s.ID()].name)

		fmt.Println("User with id " + s.ID() + " left the chat")
		fmt.Println("closed", reason)

	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
