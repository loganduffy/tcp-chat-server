package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Handling connection from: ", conn.RemoteAddr())
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
