package mapper

import (
    "grpc-exchange/generation/spot"
    "grpc-exchange/internal/domain"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func DomainMarketToProto(m *domain.Market) *spot.Market {
    protoMarket := &spot.Market{
        Id:      m.ID,
        Symbol:  m.Symbol,
        Enabled: m.Enabled,
    }
    
    if m.DeletedAt != nil {
        protoMarket.DeleteAt = timestamppb.New(*m.DeletedAt) 
				//Если объект был удален — код берет дату удаления, переводит её в специальный формат для отправки по сети и сохраняет.
    }
    
    return protoMarket
}

func ProtoMarketsToDomain(protoMarkets []*spot.Market) []*domain.Market {
    markets := make([]*domain.Market, 0, len(protoMarkets))
    for _, pm := range protoMarkets {
        market := &domain.Market{
            ID:      pm.Id,
            Symbol:  pm.Symbol,
            Enabled: pm.Enabled,
        }
        
        if pm.DeleteAt != nil && !pm.DeleteAt.AsTime().IsZero() {
            t := pm.DeleteAt.AsTime()
            market.DeletedAt = &t
        }
        
        markets = append(markets, market)
    }
    return markets
}