package mapper

import (
    "time"

    "github.com/shopspring/decimal"
    "google.golang.org/protobuf/types/known/timestamppb"

    commonv1 "grpc-exchange/gen/common"
    spotv1 "grpc-exchange/gen/spot"
    "github.com/grpc-exchange/services/spotInstrumentService/internal/domain"
)

func DecimalToProto(d decimal.Decimal) *commonv1.Decimal {
    units := d.IntPart()
    nanos := d.Sub(decimal.NewFromInt(units)).Mul(decimal.New(1, 9)).IntPart()
    return &commonv1.Decimal{
        Units: units,
        Nanos: int32(nanos),
    }
}

func MoneyToProto(d decimal.Decimal, currency string) *commonv1.Money {
    return &commonv1.Money{
        Amount:       DecimalToProto(d),
        CurrencyCode: currency,
    }
}

func TimestampToProto(t time.Time) *timestamppb.Timestamp {
    return timestamppb.New(t)
}

func TimestampPtrToProto(t *time.Time) *timestamppb.Timestamp {
    if t == nil {
        return nil
    }
    return timestamppb.New(*t)
}

func MarketToProto(m *domain.Market) *spotv1.Market {
    if m == nil {
        return nil
    }

    return &spotv1.Market{
        MarketId:     m.ID,
        Name:         m.Name,
        BaseAsset:    m.BaseAsset,
        QuoteAsset:   m.QuoteAsset,
        Enabled:      m.Enabled,
        CreatedAt:    TimestampToProto(m.CreatedAt),
        DeletedAt:    TimestampPtrToProto(m.DeletedAt),
        MinPrice:     MoneyToProto(m.MinPrice, "USD"),
        MaxPrice:     MoneyToProto(m.MaxPrice, "USD"),
        MinQuantity:  DecimalToProto(m.MinQuantity),
        MaxQuantity:  DecimalToProto(m.MaxQuantity),
    }
}