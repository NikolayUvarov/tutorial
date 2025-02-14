package main

import (
    "crypto/tls"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "google.golang.org/protobuf/proto"
    pb "main/chat"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// Handle WebSocket connection
func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket error:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }

        var req pb.ChatMessage
        if err := proto.Unmarshal(msg, &req); err != nil {
            log.Println("Protobuf decode error:", err)
            continue
        }

        log.Printf("Received from %s: %s", req.Username, req.Text)

        // Send response
        res := &pb.ChatMessage{Username: "Server", Text: "Echo: " + req.Text}
        resData, _ := proto.Marshal(res)
        conn.WriteMessage(websocket.BinaryMessage, resData)
    }
}

// Serve the client page
func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
}

func main() {
    // Load TLS certificates
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatal("Failed to load certificates:", err)
    }

    // Secure server
    server := &http.Server{
        Addr:    ":8080",
        Handler: nil,
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
        },
    }

    // Routes
    http.HandleFunc("/", serveIndex) // Serve index.html
    http.HandleFunc("/ws", handleConnection) // WebSocket handler

    // ** Serve static files (including chat.proto) **
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("WSS server running on https://localhost:8080")
    log.Fatal(server.ListenAndServeTLS("", ""))
}
