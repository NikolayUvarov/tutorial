package main

import (
    "log"
    "net"

    "google.golang.org/protobuf/proto"
    pb "main/chat"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        log.Println("Ошибка чтения:", err)
        return
    }

    var req pb.MessageRequest
    if err := proto.Unmarshal(buffer[:n], &req); err != nil {
        log.Println("Ошибка декодирования:", err)
        return
    }

    log.Printf("Сообщение от %s: %s", req.Username, req.Text)

    res := &pb.MessageResponse{Status: "Принято"}
    resData, _ := proto.Marshal(res)
    conn.Write(resData)
}

func main() {
    listener, err := net.Listen("tcp", ":5000")
    if err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
    defer listener.Close()

    log.Println("TCP-сервер запущен на порту 5000...")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Ошибка соединения:", err)
            continue
        }
        go handleConnection(conn)
    }
}
