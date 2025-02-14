package main

import (
    "crypto/tls"
    "crypto/x509"
    "log"
    "os"

    "google.golang.org/protobuf/proto"
    pb "main/chat"
)

func main() {
    // Load server certificate
    certPool := x509.NewCertPool()
    certData, err := os.ReadFile("server.crt")
    if err != nil {
        log.Fatal("Failed to read server certificate:", err)
    }
    certPool.AppendCertsFromPEM(certData)

    // Create TLS configuration
    config := &tls.Config{
        RootCAs: certPool,
    }

    // Connect to the secure server
    conn, err := tls.Dial("tcp", "localhost:50052", config)
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer conn.Close()

    // Send a Protobuf message
    msg := &pb.ChatMessage{Username: "Client", Text: "Hello, Secure Server!"}
    data, _ := proto.Marshal(msg)

    _, err = conn.Write(data)
    if err != nil {
        log.Fatal("Failed to send message:", err)
    }

    // Receive response
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        log.Fatal("Failed to read response:", err)
    }

    var res pb.ChatMessage
    err = proto.Unmarshal(buffer[:n], &res)
    if err != nil {
        log.Fatal("Failed to decode response:", err)
    }

    log.Printf("Server Response: %s", res.Text)
}
