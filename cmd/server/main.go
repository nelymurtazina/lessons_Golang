package main

import (
	"log"
	"net"
	"net/http"

	"grpc-exchange/generation/order"
	"grpc-exchange/generation/spot"
	"grpc-exchange/internal/interceptor"
	"grpc-exchange/internal/service"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Инициализация логгера zap
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync() // Сброс буферов при завершении программы

	logger.Info("Starting gRPC server...")

	spotService := service.NewSpotInstrumentService()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	spotClient := spot.NewSpotInstrumentServiceClient(conn)

	orderService := service.NewOrderService(spotClient)

	// Настраиваем цепочку интерсепторов
	// Порядок важен! Сначала восстанавливаем от паник, потом добавляем ID, потом логируем, потом метрики.
	unaryInterceptors := grpc.ChainUnaryInterceptor(
		interceptor.UnaryPanicRecoveryInterceptor(),
		interceptor.UnaryXRequestIDInterceptor(),
		interceptor.UnaryLoggerInterceptor(logger),
		interceptor.UnaryPrometheusInterceptor(),
	)

	grpcServer := grpc.NewServer(unaryInterceptors)

	spot.RegisterSpotInstrumentServiceServer(grpcServer, spotService)
	order.RegisterOrderServiceServer(grpcServer, orderService)

	// Регистрируем метрики Prometheus для gRPC
	grpc_prometheus.Register(grpcServer)

	// Запускаем HTTP сервер для отдачи метрик Prometheus (в отдельной горутине)
	go func() {
		logger.Info("Starting Prometheus metrics server on :8080")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()

	logger.Info("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}