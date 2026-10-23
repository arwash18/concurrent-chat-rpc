package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"chat-rpc/common"
)

type ChatServer struct {
	mu       sync.Mutex
	users    map[string]bool
	messages []common.Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users:    make(map[string]bool),
		messages: make([]common.Message, 0),
	}
}

func (cs *ChatServer) RegisterUser(args *common.RegisterArgs, reply *common.RegisterReply) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.users[args.Username] {
		reply.Success = false
		reply.Message = "Username already taken! Please choose another one."
		return nil
	}

	cs.users[args.Username] = true
	reply.Success = true
	reply.Message = "Joined successfully!"

	cs.messages = append(cs.messages, common.Message{
		Sender:   "SYSTEM",
		Content:  fmt.Sprintf("User '%s' has joined the chat.", args.Username),
		IsSystem: true,
	})

	fmt.Printf("[JOIN] User '%s' connected.
", args.Username)
	return nil
}

func (cs *ChatServer) SendMessage(args *common.SendMessageArgs, reply *common.Void) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.messages = append(cs.messages, common.Message{
		Sender:   args.Sender,
		Content:  args.Content,
		IsSystem: false,
	})

	fmt.Printf("[%s]: %s
", args.Sender, args.Content)
	return nil
}

func (cs *ChatServer) FetchMessages(args *common.GetMessagesArgs, reply *common.GetMessagesReply) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if args.LastIndex < len(cs.messages) {
		reply.Messages = cs.messages[args.LastIndex:]
		reply.NewIndex = len(cs.messages)
	} else {
		reply.Messages = []common.Message{}
		reply.NewIndex = args.LastIndex
	}

	return nil
}

func (cs *ChatServer) LeaveUser(args *common.RegisterArgs, reply *common.Void) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.users[args.Username] {
		delete(cs.users, args.Username)
		cs.messages = append(cs.messages, common.Message{
			Sender:   "SYSTEM",
			Content:  fmt.Sprintf("User '%s' has left the chat.", args.Username),
			IsSystem: true,
		})
		fmt.Printf("[LEAVE] User '%s' disconnected.
", args.Username)
	}
	return nil
}

func main() {
	chatServer := NewChatServer()
	err := rpc.Register(chatServer)
	if err != nil {
		fmt.Println("Error registering RPC server:", err)
		return
	}

	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat RPC Server running on port 12345...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("
Shutting down server gracefully...")
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
