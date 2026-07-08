package domain

import (
  "time"
  "github.com/shopspring/decimal"
)

// OrderSide - направление заказа (покупка/продажа)
type OrderSide string

const (
  Buy  OrderSide = "BUY"
  Sell OrderSide = "SELL"
)

// OrderType - тип заказа (лимитный/рыночный)
type OrderType string

const (
  Limit  OrderType = "LIMIT"
  Market OrderType = "MARKET"
)

// OrderStatus - статус заказа
type OrderStatus string

const (
  Pending          OrderStatus = "PENDING"
  Filled           OrderStatus = "FILLED"
  PartiallyFilled  OrderStatus = "PARTIALLY_FILLED"
  Cancelled        OrderStatus = "CANCELLED"
  Rejected         OrderStatus = "REJECTED"
)

// Order - доменная модель заказа
type Order struct {
  ID             string          // UUID
  UserID         string          // UUID пользователя
  MarketID       string          // ID рынка (например, "BTC-USD")
  Side           OrderSide       // BUY или SELL
  Type           OrderType       // LIMIT или MARKET
  Price          decimal.Decimal // Цена за единицу
  Quantity       decimal.Decimal // Количество
  FilledQuantity decimal.Decimal // Сколько исполнено
  Status         OrderStatus     // Текущий статус
  CreatedAt      time.Time
  UpdatedAt      time.Time
}