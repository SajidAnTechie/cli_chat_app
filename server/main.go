package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type user struct {
	name string
	room string
}

func main() {

	mapMake := make(map[string]*user)

	userJoin := func(userID string, joinedRoom string, userName string) *user {

		mapMake[userID] = &user{name: userName, room: joinedRoom}

		return mapMake[userID]

	}
	getJoinedUserDetails := func(userID string) *user {
		fmt.Println(userID)
		return mapMake[userID]
	}

	userLeave := func(userID string) {
		delete(mapMake, userID)
	}

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {

		fmt.Println("connected successfull")

		s.Emit("message", "Welcom to the chat")

		return nil
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomName string, userName string) {

		user := userJoin(s.ID(), roomName, userName)

		s.Join(user.room)

		server.BroadcastToRoom("/", user.room, "message", user.name+" join the chat")

	})

	server.OnEvent("/", "chatMessage", func(s socketio.Conn, msg string) {

		getJoinedUserDetails := getJoinedUserDetails(s.ID())

		server.BroadcastToRoom("/", getJoinedUserDetails.room, "message", getJoinedUserDetails.name+":  "+msg)

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

		s.Emit("leaveRoom", mapMake[s.ID()].name)

		fmt.Println("User with id " + mapMake[s.ID()].name + " left the chat")

		userLeave(s.ID())

		fmt.Println("closed", reason)

	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
