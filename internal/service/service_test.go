package service

import (
	"context"
	"testing"

	"grpc-exchange/generation/order"
	"grpc-exchange/generation/spot"
	"grpc-exchange/internal/interceptor"

	"google.golang.org/grpc"
)

// Тест для проверки ViewMarkets
func TestViewMarkets(t *testing.T) {
	// Создаем сервис
	service := NewSpotInstrumentService()

	// Создаем пустой запрос
	req := &spot.ViewMarketsRequest{}

	// Вызываем метод
	resp, err := service.ViewMarkets(context.Background(), req)

	// Проверяем что ошибки нет
	if err != nil {
		t.Errorf("Ожидалась nil ошибка, но получили: %v", err)
	}

	// Проверяем что ответ не nil
	if resp == nil {
		t.Error("Ожидался не nil ответ")
	}

	// Проверяем что рынки есть
	if len(resp.Markets) == 0 {
		t.Error("Ожидалось что рынки будут возвращены")
	}

	// Проверяем что все рынки активные
	for _, market := range resp.Markets {
		if !market.Enabled {
			t.Errorf("Рынок %s должен быть активным", market.Id)
		}
		if market.DeleteAt != nil {
			t.Errorf("Рынок %s не должен быть удален", market.Id)
		}
	}
}

// Тест для проверки CreateOrder с активным рынком
func TestCreateOrder_Success(t *testing.T) {
	// Создаем мок для spot сервиса
	mockSpot := &SimpleMockSpotService{
		ActiveMarkets: []string{"BTC_USD", "ETH_USD"},
	}

	// Создаем order сервис
	orderService := NewOrderService(mockSpot)

	// Создаем запрос
	req := &order.CreateOrderRequest{
		UserId:    "user1",
		MarketId:  "BTC_USD",
		OrderType: "buy",
		Price:     100.0,
		Quantity:  1,
	}

	// Вызываем метод
	resp, err := orderService.CreateOrder(context.Background(), req)

	// Проверяем
	if err != nil {
		t.Errorf("Не ожидалась ошибка, но получили: %v", err)
	}

	if resp == nil {
		t.Error("Ожидался ответ")
	}

	if resp.OrderId == "" {
		t.Error("Ожидался order_id")
	}

	if resp.Status != "created" {
		t.Errorf("Ожидался статус 'created', но получили: %s", resp.Status)
	}
}

// Тест для проверки CreateOrder с неактивным рынком
func TestCreateOrder_InactiveMarket(t *testing.T) {
	// Создаем мок где рынок неактивен
	mockSpot := &SimpleMockSpotService{
		ActiveMarkets: []string{}, // Пустой список - рынок неактивен
	}

	orderService := NewOrderService(mockSpot)

	req := &order.CreateOrderRequest{
		UserId:    "user1",
		MarketId:  "LTC_USD",
		OrderType: "sell",
		Price:     50.0,
		Quantity:  10,
	}

	_, err := orderService.CreateOrder(context.Background(), req)

	// Проверяем что была ошибка
	if err == nil {
		t.Error("Ожидалась ошибка для неактивного рынка")
	}
}

// Тест для проверки GetOrderStatus
func TestGetOrderStatus(t *testing.T) {
	mockSpot := &SimpleMockSpotService{
		ActiveMarkets: []string{"BTC_USD"},
	}

	orderService := NewOrderService(mockSpot)

	// Сначала создаем заказ
	createReq := &order.CreateOrderRequest{
		UserId:    "user1",
		MarketId:  "BTC_USD",
		OrderType: "buy",
		Price:     100.0,
		Quantity:  1,
	}

	createResp, _ := orderService.CreateOrder(context.Background(), createReq)

	// Теперь получаем статус
	statusReq := &order.GetOrderStatusRequest{
		OrderId: createResp.OrderId,
		UserId:  "user1",
	}

	statusResp, err := orderService.GetOrderStatus(context.Background(), statusReq)

	if err != nil {
		t.Errorf("Не ожидалась ошибка: %v", err)
	}

	if statusResp.Status != "created" {
		t.Errorf("Ожидался статус 'created', но получили: %s", statusResp.Status)
	}
}

// Тест для проверки request-id
func TestXRequestID(t *testing.T) {
	// Создаем интерсептор
	xRequestIDInterceptor := interceptor.UnaryXRequestIDInterceptor()

	// Флаг что handler был вызван
	handlerCalled := false

	// Создаем handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Получаем request-id из контекста через функцию пакета
		requestID := interceptor.GetRequestID(ctx)

		// Проверяем что он есть
		if requestID == "" {
			t.Error("request-id не должен быть пустым")
		}

		handlerCalled = true
		return "ok", nil
	}

	// Вызываем интерсептор
	ctx := context.Background()
	_, err := xRequestIDInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)

	if err != nil {
		t.Errorf("Не ожидалась ошибка: %v", err)
	}

	if !handlerCalled {
		t.Error("Handler должен быть вызван")
	}
}

// Простой мок для SpotInstrumentService
type SimpleMockSpotService struct {
	ActiveMarkets []string
}

func (m *SimpleMockSpotService) ViewMarkets(ctx context.Context, req *spot.ViewMarketsRequest, opts ...grpc.CallOption) (*spot.ViewMarketsResponse, error) {
	var markets []*spot.Market

	// Добавляем только активные рынки
	for _, marketId := range m.ActiveMarkets {
		markets = append(markets, &spot.Market{
			Id:      marketId,
			Symbol:  marketId,
			Enabled: true,
		})
	}

	return &spot.ViewMarketsResponse{
		Markets: markets,
	}, nil
}