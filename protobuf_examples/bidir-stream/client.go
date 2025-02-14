package main

import (
    "bufio"
    "context"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "main/chat"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Ошибка подключения: %v", err)
    }
    defer conn.Close()

    client := pb.NewChatServiceClient(conn)

    stream, err := client.Chat(context.Background())
    if err != nil {
        log.Fatalf("Ошибка создания стрима: %v", err)
    }

    go func() {
        for {
            res, err := stream.Recv()
            if err != nil {
                log.Fatalf("Ошибка чтения: %v", err)
            }
            log.Printf("Сервер: %s", res.Text)
        }
    }()

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()

        err := stream.Send(&pb.ChatMessage{Username: "Клиент", Text: text})
        if err != nil {
            log.Fatalf("Ошибка отправки: %v", err)
        }

        time.Sleep(time.Millisecond * 500) // Даем время серверу ответить
    }
}
