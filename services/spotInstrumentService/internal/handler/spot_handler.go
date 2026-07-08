package handler

import (
    "context"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    spotv1 "grpc-exchange/gen/spot"
    "github.com/grpc-exchange/services/spotInstrumentService/internal/mapper"
    "github.com/grpc-exchange/services/spotInstrumentService/internal/repository"
)

type SpotHandler struct {
    spotv1.UnimplementedSpotInstrumentServiceServer
    marketRepo repository.MarketRepository
}

func NewSpotHandler(marketRepo repository.MarketRepository) *SpotHandler {
    return &SpotHandler{
        marketRepo: marketRepo,
    }
}

func (h *SpotHandler) ViewMarkets(ctx context.Context, req *spotv1.ViewMarketsRequest) (*spotv1.ViewMarketsResponse, error) {
    markets, err := h.marketRepo.GetAll(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get markets: %v", err)
    }

    var activeMarkets []*spotv1.Market
    for _, m := range markets {
        // Возвращаем только активные рынки (enabled: true и deleted_at == nil)
        if m.Enabled && m.DeletedAt == nil {
            activeMarkets = append(activeMarkets, mapper.MarketToProto(m))
        }
    }

    return &spotv1.ViewMarketsResponse{
        Markets: activeMarkets,
    }, nil
}

func (h *SpotHandler) GetMarket(ctx context.Context, req *spotv1.GetMarketRequest) (*spotv1.Market, error) {
    market, err := h.marketRepo.GetByID(ctx, req.MarketId)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "market not found: %v", err)
    }
    if market == nil {
        return nil, status.Errorf(codes.NotFound, "market %s not found", req.MarketId)
    }

    return mapper.MarketToProto(market), nil
}