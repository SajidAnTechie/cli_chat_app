package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// connectCmd represents the connect command
var (
	userName   string
	roomName   string
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Command used for join rooms",

		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("userName: " + userName)
			fmt.Println("roomName: " + roomName)

			opts := &socketio_client.Options{
				Transport: "websocket",
			}

			uri := "http://localhost:8000/socket.io/"

			client, err := socketio_client.NewClient(uri, opts)
			if err != nil {
				log.Printf("NewClient error:%v\n", err)
				return
			}
			client.On("connection", func() {

				client.Emit("joinRoom", roomName, userName)

				fmt.Println("===Start Charting====")
			})

			client.On("message", func(msg string) {
				log.Printf(msg)
			})
			client.On("leaveRoom", func(userName string) {
				log.Printf(userName + " left the chat")
			})
			client.On("disconnection", func() {
				log.Printf("on disconnect\n")
			})

			client.On("error", func() {
				log.Printf("on error\n")
			})

			reader := bufio.NewReader(os.Stdin)
			for {
				data, _, _ := reader.ReadLine()
				command := string(data)
				client.Emit("chatMessage", command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.PersistentFlags().StringVar(&userName, "name", "n", "A user name")
	connectCmd.PersistentFlags().StringVar(&roomName, "room", "r", "A room name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
