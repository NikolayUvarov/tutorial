package main

import (
    "context"
    "crypto/x509"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "time"

    pb "main/example"
)

func main() {
    // Load the new server certificate
    certPool := x509.NewCertPool()
    certData, err := ioutil.ReadFile("server.crt")
    if err != nil {
        log.Fatalf("Error reading server certificate: %v", err)
    }
    if !certPool.AppendCertsFromPEM(certData) {
        log.Fatalf("Failed to add server certificate to pool")
    }

    // Create TLS credentials for secure connection
    creds := credentials.NewClientTLSFromCert(certPool, "localhost") // Use "localhost"

    // Dial the gRPC server with TLS
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewGreeterClient(conn)

    // Create a timeout context
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Call the SayHello method
    response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Golang"})
    if err != nil {
        log.Fatalf("Error calling SayHello: %v", err)
    }

    log.Printf("Server Response: %s", response.Message)
}
