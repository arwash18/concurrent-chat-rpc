# Concurrent Go RPC Chat Application

A concurrent client-server chat application in Go using standard `net/rpc` and `goroutines`.

## Features
- **Concurrent Architecture**: Supports multiple users connected simultaneously.
- **Unique Username Validation**: Prevents duplicate usernames.
- **Message Broadcasting**: Delivers messages to all other clients excluding sender.
- **Notifications**: Notifies all connected clients when someone joins or leaves.
- **Clean Disconnection**: Handles `Ctrl+C` and `/exit` signals gracefully.

## How to Run

1. **Start the Server:**
   ```bash
   go run server/main.go
   ```

2. **Start One or Multiple Clients (in separate terminal windows):**
   ```bash
   go run client/main.go
   ```
