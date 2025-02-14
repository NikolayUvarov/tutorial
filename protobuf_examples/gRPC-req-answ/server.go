package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "main/example"
)

// Реализация сервиса Greeter
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

    grpcServer := grpc.NewServer()
    pb.RegisterGreeterServer(grpcServer, &greeterServer{})

    log.Println("gRPC сервер запущен на порту 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Ошибка запуска gRPC: %v", err)
    }
}
