package main

import (
    "log"
    "net"
    "net/http"

    grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "google.golang.org/grpc"

    orderv1 "grpc-exchange/gen/order"
    "github.com/grpc-exchange/services/orderService/internal/client"
    "github.com/grpc-exchange/services/orderService/internal/handler"
    "github.com/grpc-exchange/services/orderService/internal/repository"
    "github.com/grpc-exchange/shared/interceptor"
    "github.com/grpc-exchange/shared/logger"
)

func main() {
    if err := logger.InitLogger(); err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }

    orderRepo := repository.NewInMemoryOrderRepository()

    spotClient, spotConn, err := client.NewSpotClient("localhost:50052")
    if err != nil {
        log.Fatalf("Failed to create spot client: %v", err)
    }
    defer spotConn.Close()

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            interceptor.XRequestIDUnaryInterceptor(),
            interceptor.LoggerUnaryInterceptor(),
            interceptor.PanicRecoveryUnaryInterceptor(),
            grpc_prometheus.UnaryServerInterceptor,
        ),
    )

    orderHandler := handler.NewOrderHandler(orderRepo, spotClient)
    orderv1.RegisterOrderServiceServer(s, orderHandler)

    // HTTP сервер для метрик
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Println("Metrics: http://localhost:9091/metrics")
        if err := http.ListenAndServe(":9091", nil); err != nil {
            log.Fatalf("Failed to start metrics server: %v", err)
        }
    }()

    log.Println("OrderService listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}