package domain

import (
    "time"

    "github.com/shopspring/decimal"
)

// Market - доменная модель рынка
type Market struct {
    ID          string          // ID рынка (например, "BTC-USD")
    Name        string          // Название (например, "Bitcoin/USD")
    BaseAsset   string          // Базовая валюта (например, "BTC")
    QuoteAsset  string          // Котируемая валюта (например, "USD")
    Enabled     bool            // Активен ли рынок
    CreatedAt   time.Time       // Дата создания
    DeletedAt   *time.Time      // Дата удаления (nil если активен)
    MinPrice    decimal.Decimal // Минимальная цена
    MaxPrice    decimal.Decimal // Максимальная цена
    MinQuantity decimal.Decimal // Минимальное количество
    MaxQuantity decimal.Decimal // Максимальное количество
}