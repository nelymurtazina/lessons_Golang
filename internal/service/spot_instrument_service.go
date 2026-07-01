package service

import (
    "context"
    "sync"
    "time"

    "grpc-exchange/generation/spot"
    "grpc-exchange/internal/domain"
    "grpc-exchange/internal/mapper"
)

type SpotInstrumentService struct {
    spot.UnimplementedSpotInstrumentServiceServer
    
    markets map[string]*domain.Market 
    mu      sync.RWMutex
}

func NewSpotInstrumentService() *SpotInstrumentService {
    svc := &SpotInstrumentService{
        markets: make(map[string]*domain.Market),
    }
    
    // Инициализация тестовыми данными
    now := time.Now()
    
    svc.markets["BTC_USD"] = &domain.Market{
        ID:        "BTC_USD",
        Symbol:    "BTC/USD",
        Enabled:   true,
        DeletedAt: nil, 
    }
    
    svc.markets["ETH_USD"] = &domain.Market{
        ID:        "ETH_USD",
        Symbol:    "ETH/USD",
        Enabled:   true,
        DeletedAt: nil,
    }
    
    svc.markets["LTC_USD"] = &domain.Market{
        ID:        "LTC_USD",
        Symbol:    "LTC/USD",
        Enabled:   false, // Приостановлен
        DeletedAt: nil,
    }
    
    deletedTime := now.Add(-24 * time.Hour) 
    svc.markets["XRP_USD"] = &domain.Market{
        ID:        "XRP_USD",
        Symbol:    "XRP/USD",
        Enabled:   true,
        DeletedAt: &deletedTime, 
    }
    
    return svc
}

// ViewMarkets возвращает список доступных рынков
// Требование: только enabled=true и deleted_at=null
func (s *SpotInstrumentService) ViewMarkets(
    ctx context.Context, 
    req *spot.ViewMarketsRequest,
) (*spot.ViewMarketsResponse, error) {
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var activeMarkets []*spot.Market
    
    for _, market := range s.markets {
        if market.Enabled && market.DeletedAt == nil {
            protoMarket := mapper.DomainMarketToProto(market)
            activeMarkets = append(activeMarkets, protoMarket)
        }
    }
    
    return &spot.ViewMarketsResponse{
        Markets: activeMarkets,
    }, nil
}

// GetMarket - вспомогательный метод для получения конкретного рынка
// Используется OrderService для проверки перед созданием заказа
func (s *SpotInstrumentService) GetMarket(marketID string) (*domain.Market, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    market, exists := s.markets[marketID]
    if !exists {
        return nil, ErrMarketNotFound
    }
    
    if !market.Enabled {
        return nil, ErrMarketNotActive
    }
    
    if market.DeletedAt != nil {
        return nil, ErrMarketNotActive
    }
    
    return market, nil
}