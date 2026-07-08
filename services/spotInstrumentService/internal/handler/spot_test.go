package handler_test

import (
	"context"
	"testing"

	spotv1 "grpc-exchange/gen/spot"

	"github.com/grpc-exchange/services/spotInstrumentService/internal/handler"
	"github.com/grpc-exchange/services/spotInstrumentService/internal/repository"
)

func TestViewMarkets_ReturnsOnlyActiveMarkets(t *testing.T) {
    //  ПОДГОТОВКА
    marketRepo := repository.NewInMemoryMarketRepository()
    spotHandler := handler.NewSpotHandler(marketRepo)

    ctx := context.Background()
    req := &spotv1.ViewMarketsRequest{
        UserRoles: []string{"user"},
    }

    // ВЫПОЛНЕНИЕ
    resp, err := spotHandler.ViewMarkets(ctx, req)

    // ПРОВЕРКА
    if err != nil {
        t.Fatalf("ViewMarkets failed: %v", err)
    }

    // Проверяем, что возвращены только активные рынки
    for _, market := range resp.Markets {
        if !market.Enabled {
            t.Errorf("Market %s is not enabled but returned", market.MarketId)
        }
        if market.DeletedAt != nil {
            t.Errorf("Market %s is deleted but returned", market.MarketId)
        }
    }

    // Проверяем количество активных рынков 
    expectedCount := 2
    if len(resp.Markets) != expectedCount {
        t.Errorf("Expected %d active markets, got %d", expectedCount, len(resp.Markets))
    }
}

func TestViewMarkets_WithEmptyUserRoles(t *testing.T) {
    marketRepo := repository.NewInMemoryMarketRepository()
    spotHandler := handler.NewSpotHandler(marketRepo)

    ctx := context.Background()
    req := &spotv1.ViewMarketsRequest{
        UserRoles: []string{},
    }

    resp, err := spotHandler.ViewMarkets(ctx, req)
    if err != nil {
        t.Fatalf("ViewMarkets failed: %v", err)
    }

    if len(resp.Markets) == 0 {
        t.Error("Expected markets with empty roles, got 0")
    }
}