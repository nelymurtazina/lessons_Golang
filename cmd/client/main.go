package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-exchange/generation/order"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	clientInterceptor := grpc.WithUnaryInterceptor(
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			requestID := uuid.New().String()

			// Добавляем метаданные
			md := metadata.Pairs("request-id", requestID)
			ctx = metadata.NewOutgoingContext(ctx, md)

			log.Printf("[Client] Sending request with request-id: %s", requestID)

			return invoker(ctx, method, req, reply, cc, opts...)
		},
	)

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Без TLS, вроде норм
		clientInterceptor,
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := order.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Testing CreateOrder")

	createReq := &order.CreateOrderRequest{
		UserId:    "user_123",
		MarketId:  "BTC_USD", // Активный рынок
		OrderType: "buy",
		Price:     50000.0,
		Quantity:  1,
		UserRoles: []string{"trader"},
	}

	createResp, err := client.CreateOrder(ctx, createReq)
	if err != nil {
		log.Printf("Error creating order: %v", err)
	} else {
		fmt.Printf("Order created successfully!\n")
		fmt.Printf("Order ID: %s\n", createResp.OrderId)
		fmt.Printf("Status: %s\n", createResp.Status)
	}

	// Тесли заказ создан
	if err == nil && createResp.OrderId != "" {
		fmt.Println("\n=== Testing GetOrderStatus ===")

		statusReq := &order.GetOrderStatusRequest{
			OrderId: createResp.OrderId,
			UserId:  "user_123",
		}

		statusResp, err := client.GetOrderStatus(ctx, statusReq)
		if err != nil {
			log.Printf("Error getting order status: %v", err)
		} else {
			fmt.Printf("Order status retrieved!\n")
			fmt.Printf("Order ID: %s\n", createResp.OrderId)
			fmt.Printf("Status: %s\n", statusResp.Status)
		}
	}

	fmt.Println("\n Testing CreateOrder with inactive market (LTC_USD)")

	createReqInactive := &order.CreateOrderRequest{
		UserId:    "user_456",
		MarketId:  "LTC_USD", // Неактивный рынок
		OrderType: "sell",
		Price:     100.0,
		Quantity:  10,
		UserRoles: []string{"trader"},
	}

	_, err = client.CreateOrder(ctx, createReqInactive)
	if err != nil {
		fmt.Printf("Expected error received: %v\n", err)
	} else {
		fmt.Printf("Should have failed for inactive market!\n")
	}

	fmt.Println("\n Tests completed")
}
