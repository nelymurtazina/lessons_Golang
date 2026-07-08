package main

import (
    "log"
    "net"
    "net/http"

    "github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "google.golang.org/grpc"

    spotv1 "grpc-exchange/gen/spot"
    "github.com/grpc-exchange/services/spotInstrumentService/internal/handler"
    "github.com/grpc-exchange/services/spotInstrumentService/internal/repository"
    "github.com/grpc-exchange/shared/interceptor"
    "github.com/grpc-exchange/shared/logger"
)

func main() {
    if err := logger.InitLogger(); err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }

    marketRepo := repository.NewInMemoryMarketRepository()

    lis, err := net.Listen("tcp", ":50052")
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

    spotHandler := handler.NewSpotHandler(marketRepo)
    spotv1.RegisterSpotInstrumentServiceServer(s, spotHandler)

    // HTTP сервер для метрик
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Println("Metrics: http://localhost:9090/metrics")
        if err := http.ListenAndServe(":9090", nil); err != nil {
            log.Fatalf("Failed to start metrics server: %v", err)
        }
    }()

    log.Println("SpotInstrumentService listening on :50052")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}