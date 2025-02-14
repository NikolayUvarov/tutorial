package main

import (
    "io"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "main/chat"
)

type chatServer struct {
    pb.UnimplementedChatServiceServer
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
    for {
        msg, err := stream.Recv() // Получаем сообщение от клиента
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        log.Printf("Сообщение от %s: %s", msg.Username, msg.Text)

        // Отправляем ответ обратно
        err = stream.Send(&pb.ChatMessage{
            Username: "Сервер",
            Text:     "Эхо: " + msg.Text,
        })
        if err != nil {
            return err
        }
    }
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterChatServiceServer(grpcServer, &chatServer{})

    log.Println("gRPC сервер запущен на порту 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Ошибка запуска gRPC: %v", err)
    }
}
