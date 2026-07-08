package repository

import (
  "context"
  "sync"

  "github.com/shopspring/decimal"
  "github.com/grpc-exchange/services/spotInstrumentService/internal/domain"
)

type MarketRepository interface {
    GetAll(ctx context.Context) ([]*domain.Market, error)
    GetByID(ctx context.Context, id string) (*domain.Market, error)
}

type InMemoryMarketRepository struct {
    mu      sync.RWMutex
    markets map[string]*domain.Market
}

func NewInMemoryMarketRepository() *InMemoryMarketRepository {
    repo := &InMemoryMarketRepository{
        markets: make(map[string]*domain.Market),
    }

    // Добавляем тестовые рынки
    repo.markets["BTC-USD"] = &domain.Market{
        ID:          "BTC-USD",
        Name:        "Bitcoin/USD",
        BaseAsset:   "BTC",
        QuoteAsset:  "USD",
        Enabled:     true,
        DeletedAt:   nil,
        MinPrice:    decimal.NewFromFloat(10000),
        MaxPrice:    decimal.NewFromFloat(100000),
        MinQuantity: decimal.NewFromFloat(0.001),
        MaxQuantity: decimal.NewFromFloat(100),
    }
    repo.markets["ETH-USD"] = &domain.Market{
        ID:          "ETH-USD",
        Name:        "Ethereum/USD",
        BaseAsset:   "ETH",
        QuoteAsset:  "USD",
        Enabled:     true,
        DeletedAt:   nil,
        MinPrice:    decimal.NewFromFloat(500),
        MaxPrice:    decimal.NewFromFloat(5000),
        MinQuantity: decimal.NewFromFloat(0.01),
        MaxQuantity: decimal.NewFromFloat(1000),
    }

    return repo
}

func (r *InMemoryMarketRepository) GetAll(ctx context.Context) ([]*domain.Market, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    markets := make([]*domain.Market, 0, len(r.markets))
    for _, m := range r.markets {
        markets = append(markets, m)
    }
    return markets, nil
}

func (r *InMemoryMarketRepository) GetByID(ctx context.Context, id string) (*domain.Market, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.markets[id], nil
}