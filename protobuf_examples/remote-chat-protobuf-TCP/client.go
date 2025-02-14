package main

import (
    "log"
    "net"

    "google.golang.org/protobuf/proto"
    pb "main/chat"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:5000")
    if err != nil {
        log.Fatal("Ошибка подключения:", err)
    }
    defer conn.Close()

    req := &pb.MessageRequest{Username: "Alice", Text: "Привет, мир!"}
    data, _ := proto.Marshal(req)

    conn.Write(data)

    buffer := make([]byte, 1024)
    n, _ := conn.Read(buffer)

    var res pb.MessageResponse
    proto.Unmarshal(buffer[:n], &res)

    log.Println("Ответ сервера:", res.Status)
}
