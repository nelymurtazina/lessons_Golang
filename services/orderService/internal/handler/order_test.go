package handler_test

import (
    "context"
    "testing"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    commonv1 "grpc-exchange/gen/common"
    orderv1 "grpc-exchange/gen/order"
    spotv1 "grpc-exchange/gen/spot"
    "github.com/grpc-exchange/services/orderService/internal/handler"
    "github.com/grpc-exchange/services/orderService/internal/repository"
)

type mockSpotClient struct {
    markets []*spotv1.Market
    err     error
}

func (m *mockSpotClient) ViewMarkets(ctx context.Context, req *spotv1.ViewMarketsRequest, opts ...grpc.CallOption) (*spotv1.ViewMarketsResponse, error) {
    if m.err != nil {
        return nil, m.err
    }
    return &spotv1.ViewMarketsResponse{
        Markets: m.markets,
    }, nil
}

func (m *mockSpotClient) GetMarket(ctx context.Context, req *spotv1.GetMarketRequest, opts ...grpc.CallOption) (*spotv1.Market, error) {
    return nil, nil
}

func TestCreateOrder_RejectsInvalidMarket(t *testing.T) {
    orderRepo := repository.NewInMemoryOrderRepository()

    mockSpot := &mockSpotClient{
        markets: []*spotv1.Market{
            {MarketId: "BTC-USD", Enabled: true},
        },
    }

    orderHandler := handler.NewOrderHandler(orderRepo, mockSpot)

    ctx := context.Background()
    req := &orderv1.CreateOrderRequest{
        UserId:   "user-123",
        MarketId: "NOT-EXISTS", 
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

    resp, err := orderHandler.CreateOrder(ctx, req)

    if err == nil {
        t.Error("Expected error for invalid market, got nil")
    }
    if resp != nil {
        t.Error("Expected nil response for invalid market")
    }

    st, ok := status.FromError(err)
    if !ok {
        t.Fatalf("Expected gRPC status error, got %v", err)
    }
    if st.Code() != codes.NotFound {
        t.Errorf("Expected NotFound code, got %s", st.Code())
    }
}

func TestCreateOrder_AcceptsValidMarket(t *testing.T) {
    orderRepo := repository.NewInMemoryOrderRepository()

    mockSpot := &mockSpotClient{
        markets: []*spotv1.Market{
            {MarketId: "BTC-USD", Enabled: true},
        },
    }

    orderHandler := handler.NewOrderHandler(orderRepo, mockSpot)

    ctx := context.Background()
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

    resp, err := orderHandler.CreateOrder(ctx, req)

    // ПРОВЕРКА
    if err != nil {
        t.Fatalf("CreateOrder failed: %v", err)
    }

    if resp.OrderId == "" {
        t.Error("Expected order ID, got empty")
    }
    if resp.Status != orderv1.OrderStatus_ORDER_STATUS_PENDING {
        t.Errorf("Expected PENDING status, got %s", resp.Status)
    }
}