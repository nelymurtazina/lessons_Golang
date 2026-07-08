package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    commonv1 "grpc-exchange/gen/common"
    orderv1 "grpc-exchange/gen/order"
)

func main() {
    log.Println("Starting Order Client...")
    // Подключаемся к OrderService (БЕЗ интерсептора)
    conn, err := grpc.NewClient(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := orderv1.NewOrderServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    log.Println("Creating order...")

    req := &orderv1.CreateOrderRequest{
      UserId:   "user-123",
      MarketId: "BTC-USD",
      Side:     orderv1.OrderSide_ORDER_SIDE_BUY,
      Type:     orderv1.OrderType_ORDER_TYPE_LIMIT,
	      Price: &commonv1.Money{
        Amount: &commonv1.Decimal{
        Units: 50000,
        Nanos: 0,
     	 },
      	CurrencyCode: "USD",
      },
        Quantity: &commonv1.Decimal{
          Units: 1,
          Nanos: 0,
      },
    }

    resp, err := client.CreateOrder(ctx, req)
    if err != nil {
        log.Fatalf("CreateOrder failed: %v", err)
    }

    log.Printf("Order created!")
    log.Printf("ID:     %s", resp.OrderId)
    log.Printf("Status: %s", resp.Status.String())

    log.Println("Getting order status...")

    statusReq := &orderv1.GetOrderStatusRequest{
        OrderId: resp.OrderId,
        UserId:  "user-123",
    }

    statusResp, err := client.GetOrderStatus(ctx, statusReq)
    if err != nil {
        log.Fatalf("GetOrderStatus failed: %v", err)
    }

    log.Printf("Order status:")
    log.Printf("ID:     %s", statusResp.Order.OrderId)
    log.Printf("Status: %s", statusResp.Order.Status.String())
    log.Printf("Market: %s", statusResp.Order.MarketId)

    log.Println("Done!")
}