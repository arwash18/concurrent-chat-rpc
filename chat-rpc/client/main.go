package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"chat-rpc/common"
)

var client *rpc.Client
var username string
var lastIndex int = 0

func pollMessages() {
	for {
		time.Sleep(500 * time.Millisecond)

		args := common.GetMessagesArgs{Username: username, LastIndex: lastIndex}
		var reply common.GetMessagesReply

		err := client.Call("ChatServer.FetchMessages", &args, &reply)
		if err != nil {
			return
		}

		for _, msg := range reply.Messages {
			if msg.Sender != username {
				if msg.IsSystem {
					fmt.Printf("
📢 [%s]
> ", msg.Content)
				} else {
					fmt.Printf("
💬 [%s]: %s
> ", msg.Sender, msg.Content)
				}
			}
		}
		lastIndex = reply.NewIndex
	}
}

func main() {
	var err error
	client, err = rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Println("Unable to connect to server! Make sure server is running on port 12345.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter unique username: ")
		username, _ = reader.ReadString('
')
		username = strings.TrimSpace(username)

		if username == "" {
			continue
		}

		var reply common.RegisterReply
		err := client.Call("ChatServer.RegisterUser", &common.RegisterArgs{Username: username}, &reply)
		if err == nil && reply.Success {
			fmt.Printf("Welcome %s! You are now connected to the chat.
", username)
			break
		}
		fmt.Println("Error:", reply.Message)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		var void common.Void
		client.Call("ChatServer.LeaveUser", &common.RegisterArgs{Username: username}, &void)
		fmt.Println("
Disconnected cleanly from server. Goodbye!")
		os.Exit(0)
	}()

	go pollMessages()

	fmt.Println("Type your message and press Enter (Type '/exit' to leave):")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('
')
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		if text == "/exit" {
			var void common.Void
			client.Call("ChatServer.LeaveUser", &common.RegisterArgs{Username: username}, &void)
			fmt.Println("Left the chat.")
			break
		}

		var void common.Void
		client.Call("ChatServer.SendMessage", &common.SendMessageArgs{Sender: username, Content: text}, &void)
	}
}
