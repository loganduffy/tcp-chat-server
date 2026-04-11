# TCP Chat Server

This is a TCP Chat Server that uses connects x amount of clients to the server.

## What I Learned

I learned that a Tranmission Control Protocol (TCP) is an exact 1 to 1 connection between a client and a server. I learned how goroutines work and security concerns involving them such as race conditions, when multiple goroutines access shared package-level data simultaneously. I handled this issue by using a mutex to lock the shared data (the clients slice) during reads and writes.

## How It Works

The server opens with a TCP listener on the port 8080 (technically could be 1024-65535) and waits for connections. When a client connects it spawns a goroutine so that multiple clients can connect at once without it hanging on the one client. Each client then gets added to a shared slice. When the client sends a message the server receives the message and then broadcasts it to all other connected clients in the slice. When a client disconnects they are simply removed from the slice which runs in O(n) time. All access to the clients slice is protected with a mutex to prevent race conditions.

## How To Run

1. Clone the repository with HTTPS `https://github.com/loganduffy/tcp-server-chat.git` 
```
git clone https://github.com/loganduffy/tcp-server-chat.git
```

2. Then from the project directory run
```
go run .
```

3. Then start connecting with netcat or telnet

```
nc localhost 8080
# or
telnet localhost 8080
```

4. Start connecting with more terminals!
