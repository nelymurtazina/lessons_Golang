package service

import (
    "context"
    "sync"
    "time"

    "github.com/google/uuid"
    "grpc-exchange/generation/order"
    "grpc-exchange/generation/spot"
    "grpc-exchange/internal/domain"
)

// OrderService реализует интерфейс OrderServiceServer
type OrderService struct {
    order.UnimplementedOrderServiceServer
    
    orders      map[string]*domain.Order 
    mu          sync.RWMutex
    spotClient  spot.SpotInstrumentServiceClient
}

func NewOrderService(spotClient spot.SpotInstrumentServiceClient) *OrderService {
    return &OrderService{
        orders:     make(map[string]*domain.Order),
        spotClient: spotClient,
    }
}

func (s *OrderService) SetSpotClient(client spot.SpotInstrumentServiceClient) {
    s.spotClient = client
}

func (s *OrderService) CreateOrder(
    ctx context.Context, 
    req *order.CreateOrderRequest,
) (*order.CreateOrderResponse, error) {
    
    if req.UserId == "" {
        return nil, ErrInvalidArgument
    }
    if req.MarketId == "" {
        return nil, ErrInvalidArgument
    }
    if req.OrderType == "" {
        return nil, ErrInvalidArgument
    }
    if req.Price <= 0 {
        return nil, ErrInvalidArgument
    }
    if req.Quantity <= 0 {
        return nil, ErrInvalidArgument
    }
    
    if s.spotClient != nil {
        marketsReq := &spot.ViewMarketsRequest{
            UserRoles: req.UserRoles,
        }
        
        marketsResp, err := s.spotClient.ViewMarkets(ctx, marketsReq)
        if err != nil {
            return nil, err
        }
        
        marketFound := false
        for _, market := range marketsResp.Markets {
            if market.Id == req.MarketId && market.Enabled {
                marketFound = true
                break
            }
        }
        
        if !marketFound {
            return nil, ErrMarketNotFound
        }
    }
    
    orderID := uuid.New().String() 
    
    newOrder := &domain.Order{
        ID:        orderID,
        UserID:    req.UserId,
        MarketID:  req.MarketId,
        OrderType: req.OrderType,
        Price:     req.Price,
        Quantity:  req.Quantity,
        Status:    "created", // Начальный статус
        CreatedAt: time.Now(),
    }
    
    s.mu.Lock()
    s.orders[orderID] = newOrder
    s.mu.Unlock()
    
    return &order.CreateOrderResponse{
        OrderId: orderID,
        Status:  newOrder.Status,
    }, nil
}

func (s *OrderService) GetOrderStatus(
    ctx context.Context, 
    req *order.GetOrderStatusRequest,
) (*order.GetOrderStatusResponse, error) {
    
    if req.OrderId == "" {
        return nil, ErrInvalidArgument
    }
    if req.UserId == "" {
        return nil, ErrInvalidArgument
    }
    
    s.mu.RLock()
    ord, exists := s.orders[req.OrderId]
    s.mu.RUnlock()
    
    if !exists {
        return nil, ErrOrderNotFound
    }
    
    // Проверка прав доступа
    // Требование: пользователь может смотреть только свои заказы
    if ord.UserID != req.UserId {
        return nil, ErrPermissionDenied
    }
    
    return &order.GetOrderStatusResponse{
        Status: ord.Status,
    }, nil
}