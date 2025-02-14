package main

import (
    "context"
    "crypto/tls"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    pb "main/example"
)

// Реализация сервиса
type greeterServer struct {
    pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    message := fmt.Sprintf("Привет, %s!", req.Name)
    return &pb.HelloResponse{Message: message}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }

    // Загружаем сертификаты
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatalf("Ошибка загрузки сертификатов: %v", err)
    }

    creds := credentials.NewServerTLSFromCert(&cert)
    grpcServer := grpc.NewServer(grpc.Creds(creds))

    pb.RegisterGreeterServer(grpcServer, &greeterServer{})

    log.Println("gRPC сервер с TLS запущен на порту 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Ошибка запуска gRPC: %v", err)
    }
}
