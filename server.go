package main

import (
	"fmt"
	"net"
	"bufio"
	"sync"
)

var clients []net.Conn
var mutex sync.Mutex

func broadcast(message string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	mutex.Lock()
	clients = append(clients, conn)
	mutex.Unlock()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			mutex.Lock()
			for i, client := range clients {
				if client == conn {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			mutex.Unlock()
			return
		}
		fmt.Println("Message received:", message)
		broadcast(message, conn)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is open and listening on port 8080")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting the connection", err)
			continue
		}

		// Go routines allowing for parallel threading within the chatroom
		go handleConnection(conn)
	}
}
