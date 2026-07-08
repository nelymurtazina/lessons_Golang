package handler

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    commonv1 "grpc-exchange/gen/common"
    orderv1 "grpc-exchange/gen/order"
    spotv1 "grpc-exchange/gen/spot"

    "github.com/grpc-exchange/services/orderService/internal/domain"
    "github.com/grpc-exchange/services/orderService/internal/repository"
)

type OrderHandler struct {
	orderv1.UnimplementedOrderServiceServer
	orderRepo  repository.OrderRepository
	spotClient spotv1.SpotInstrumentServiceClient
}

func NewOrderHandler(orderRepo repository.OrderRepository, spotClient spotv1.SpotInstrumentServiceClient) *OrderHandler {
	return &OrderHandler{
	orderRepo:  orderRepo,
	spotClient: spotClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	// Проверяем рынок через SpotInstrumentService
	if req.UserId == "" {
    return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
  }
  if req.MarketId == "" {
    return nil, status.Errorf(codes.InvalidArgument, "market_id is required")
  }
  if req.Price == nil || req.Price.Amount == nil {
    return nil, status.Errorf(codes.InvalidArgument, "price is required")
  }
  if req.Price.Amount.Units <= 0 && req.Price.Amount.Nanos <= 0 {
    return nil, status.Errorf(codes.InvalidArgument, "price must be greater than 0")
  }
  if req.Quantity == nil {
    return nil, status.Errorf(codes.InvalidArgument, "quantity is required")
  }
  if req.Quantity.Units <= 0 && req.Quantity.Nanos <= 0 {
    return nil, status.Errorf(codes.InvalidArgument, "quantity must be greater than 0")
	}

	marketsResp, err := h.spotClient.ViewMarkets(ctx, &spotv1.ViewMarketsRequest{
		UserRoles: []string{"user"},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check market: %v", err)
	}

  // Ищем market_id в списке активных рынков
	marketFound := false
	for _, m := range marketsResp.Markets {
		if m.MarketId == req.MarketId {
			marketFound = true
			break
		}
	}
	if !marketFound {
		return nil, status.Errorf(codes.NotFound, "market %s not found or inactive", req.MarketId)
	}

	// Конвертируем proto → decimal
	price := decimal.NewFromInt(req.Price.Amount.Units).Add(decimal.NewFromInt(int64(req.Price.Amount.Nanos)).Div(decimal.New(1, 9)))
	quantity := decimal.NewFromInt(req.Quantity.Units).Add(decimal.NewFromInt(int64(req.Quantity.Nanos)).Div(decimal.New(1, 9)))

	// Создаем доменный заказ
	order := &domain.Order{
		ID:             uuid.New().String(),
    UserID:         req.UserId,
    MarketID:       req.MarketId,
    Side:           domain.OrderSide(req.Side.String()),
    Type:           domain.OrderType(req.Type.String()),
    Price:          price,
    Quantity:       quantity,
    FilledQuantity: decimal.Zero,
    Status:         domain.Pending,
    CreatedAt:      time.Now(),
    UpdatedAt:      time.Now(),
  }

	// Сохраняем
	if err := h.orderRepo.Create(ctx, order); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save order: %v", err)
	}

	return &orderv1.CreateOrderResponse{
		OrderId: order.ID,
		Status:  orderv1.OrderStatus_ORDER_STATUS_PENDING,
	}, nil
}

func (h *OrderHandler) GetOrderStatus(ctx context.Context, req *orderv1.GetOrderStatusRequest) (*orderv1.GetOrderStatusResponse, error) {
	order, err := h.orderRepo.GetByID(ctx, req.OrderId)
	if err != nil {
		 return nil, status.Errorf(codes.NotFound, "order not found")
	}
	if order == nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	if order.UserID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	return &orderv1.GetOrderStatusResponse{
		Order: &orderv1.Order{
			OrderId: order.ID,
			UserId:  order.UserID,
			MarketId: order.MarketID,
			Side:    orderv1.OrderSide(orderv1.OrderSide_value[string(order.Side)]),
			Type:    orderv1.OrderType(orderv1.OrderType_value[string(order.Type)]),
			Price: &commonv1.Money{
				Amount: &commonv1.Decimal{
					Units: order.Price.IntPart(),
					Nanos: int32(order.Price.Sub(decimal.NewFromInt(order.Price.IntPart())).Mul(decimal.New(1, 9)).IntPart()),
				},
				CurrencyCode: "USD",
			},
			Quantity: &commonv1.Decimal{
				Units: order.Quantity.IntPart(),
				Nanos: int32(order.Quantity.Sub(decimal.NewFromInt(order.Quantity.IntPart())).Mul(decimal.New(1, 9)).IntPart()),
			},
			FilledQuantity: &commonv1.Decimal{
			Units: order.FilledQuantity.IntPart(),
			Nanos: int32(order.FilledQuantity.Sub(decimal.NewFromInt(order.FilledQuantity.IntPart())).Mul(decimal.New(1, 9)).IntPart()),
			},
			Status: orderv1.OrderStatus(orderv1.OrderStatus_value[string(order.Status)]),
			
		},
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *orderv1.CancelOrderRequest) (*orderv1.CancelOrderResponse, error) {
  if req.OrderId == "" {
    return nil, status.Errorf(codes.InvalidArgument, "order_id is required")
  }
  if req.UserId == "" {
    return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
  }

  order, err := h.orderRepo.GetByID(ctx, req.OrderId)
  if err != nil {
    return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
  }
  if order == nil {
    return nil, status.Errorf(codes.NotFound, "order %s not found", req.OrderId)
  }

	if order.UserID != req.UserId {
    return nil, status.Errorf(codes.PermissionDenied, "access denied")
  }

  // Можно отменить только PENDING заказ
  if order.Status != domain.Pending {
    return nil, status.Errorf(codes.FailedPrecondition, "order %s cannot be cancelled, current status: %s", req.OrderId, order.Status)
  }

  order.Status = domain.Cancelled
  order.UpdatedAt = time.Now()

  if err := h.orderRepo.Update(ctx, order); err != nil {
    return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
  }

  return &orderv1.CancelOrderResponse{
    OrderId: req.OrderId,
    Status:  orderv1.OrderStatus_ORDER_STATUS_CANCELLED,
    Message: "Order cancelled successfully",
  }, nil
}