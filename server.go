package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn     net.Conn
	username string
}

var (
	// Map usernames to Client objects for O(1) private messaging lookups
	clients = make(map[string]*Client)
	mutex   sync.Mutex
)

// Broadcast sends a message to everyone
func broadcast(message string, excludeUser string) {
	mutex.Lock()
	defer mutex.Unlock()

	for username, client := range clients {
		if username != excludeUser {
			client.conn.Write([]byte(message + "\n"))
		}
	}
}

// SendPrivateMessage routes a message directly to a specific target user
func sendPrivateMessage(sender string, target string, text string) {
	mutex.Lock()
	targetClient, exists := clients[target]
	senderClient := clients[sender]
	mutex.Unlock()

	if !exists {
		senderClient.conn.Write([]byte(fmt.Sprintf("[System] User '%s' not found.\n", target)))
		return
	}

	pm := fmt.Sprintf("[PM from %s]: %s\n", sender, text)
	targetClient.conn.Write([]byte(pm))
	senderClient.conn.Write([]byte(fmt.Sprintf("[PM to %s]: %s\n", target, text)))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 1. Prompt and registration phase
	conn.Write([]byte("Enter your username: "))
	username, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	username = strings.TrimSpace(username)

	if username == "" {
		conn.Write([]byte("Invalid username. Disconnecting.\n"))
		return
	}

	client := &Client{conn: conn, username: username}

	mutex.Lock()
	// Handle duplicate username edge case
	if _, exists := clients[username]; exists {
		mutex.Unlock()
		conn.Write([]byte("Username already taken. Disconnecting.\n"))
		return
	}
	clients[username] = client
	mutex.Unlock()

	// 2. Announce Join Event
	fmt.Printf("%s joined the chat from %s\n", username, conn.RemoteAddr())
	broadcast(fmt.Sprintf("*** %s has joined the chat ***", username), username)

	// Clean up and announce Leave Event on disconnect
	defer func() {
		mutex.Lock()
		delete(clients, username)
		mutex.Unlock()

		fmt.Printf("%s disconnected\n", username)
		broadcast(fmt.Sprintf("*** %s has left the chat ***", username), "")
	}()

	// 3. Command Parsing Loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Protocol Parsing: Check if input starts with a command
		if strings.HasPrefix(message, "/msg ") {
			// Format: /msg <username> <message text>
			parts := strings.SplitN(message, " ", 3)
			if len(parts) < 3 {
				conn.Write([]byte("[System] Usage: /msg <username> <message>\n"))
				continue
			}

			targetUser := parts[1]
			pmText := parts[2]
			sendPrivateMessage(username, targetUser, pmText)

		} else if message == "/list" {
			// Optional command: list connected users
			mutex.Lock()
			var activeUsers []string
			for u := range clients {
				activeUsers = append(activeUsers, u)
			}
			mutex.Unlock()
			conn.Write([]byte(fmt.Sprintf("[System] Online users: %s\n", strings.Join(activeUsers, ", "))))

		} else {
			// Default behavior: Broadcast to all other clients
			formattedMsg := fmt.Sprintf("[%s]: %s", username, message)
			broadcast(formattedMsg, username)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}