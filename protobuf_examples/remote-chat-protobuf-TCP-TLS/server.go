package main

import (
    "crypto/tls"
    "log"
    "net"

    "google.golang.org/protobuf/proto"
    pb "main/chat"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    for {
        // Read incoming message
        buffer := make([]byte, 1024)
        n, err := conn.Read(buffer)
        if err != nil {
            log.Println("Connection closed:", err)
            return
        }

        // Decode Protobuf message
        var req pb.ChatMessage
        if err := proto.Unmarshal(buffer[:n], &req); err != nil {
            log.Println("Protobuf decode error:", err)
            continue
        }

        log.Printf("Received from %s: %s", req.Username, req.Text)

        // Create response message
        res := &pb.ChatMessage{Username: "Server", Text: "Echo: " + req.Text}
        resData, _ := proto.Marshal(res)

        // Send response over TLS
        _, err = conn.Write(resData)
        if err != nil {
            log.Println("Write error:", err)
            return
        }
    }
}

func main() {
    // Load TLS certificates
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatal("Failed to load TLS certificates:", err)
    }

    // Create TLS listener
    config := &tls.Config{Certificates: []tls.Certificate{cert}}
    listener, err := tls.Listen("tcp", ":50052", config)
    if err != nil {
        log.Fatal("Error starting TLS server:", err)
    }
    defer listener.Close()

    log.Println("Secure TCP server running on port 50052...")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Connection error:", err)
            continue
        }
        go handleConnection(conn) // Handle each connection in a goroutine
    }
}
