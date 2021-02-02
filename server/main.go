package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type user struct {
	id       string
	roomName string
}

func main() {

	var users []user

	userJoin := func(userID string, joinedRoom string) {

		results := []user{user{
			id:       userID,
			roomName: joinedRoom,
		}}

		for _, details := range results {
			users = append(users, user{
				id:       details.id,
				roomName: details.roomName,
			})
		}

	}
	getUserJoinedRoom := func(userID string) string {

		var dat []map[string]string

		if err := json.Unmarshal(users, &dat); err != nil {
			panic(err)
		}

		for idx := range dat {
			if dat[idx]["id"] == userID {
				return dat[idx]["roomName"]
			}
		}
	}

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {

		fmt.Println("connected successfull")

		s.Emit("message", "Welcom to the chat")

		return nil
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomName string) {

		// s.Join(roomName)
		userJoin(s.ID(), roomName)

		server.JoinRoom("/", roomName, s)

		server.BroadcastToRoom("/", roomName, "message", "User with id "+s.ID()+" join the chat")

	})

	server.OnEvent("/", "chatMessage", func(s socketio.Conn, msg string) {

		getJoinedRoom := getUserJoinedRoom(s.ID())

		server.BroadcastToRoom("/", getJoinedRoom, "message", msg)

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

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
