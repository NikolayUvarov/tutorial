package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "main/example"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Ошибка подключения: %v", err)
    }
    defer conn.Close()

    client := pb.NewGreeterClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Голанг"})
    if err != nil {
        log.Fatalf("Ошибка вызова метода: %v", err)
    }

    log.Printf("Ответ сервера: %s", response.Message)
}
