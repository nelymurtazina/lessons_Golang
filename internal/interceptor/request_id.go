package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)


type contextKey string

const RequestIDKey contextKey = "request-id"

//UnaryServerInterceptor-Это middleware, которая срабатывает на сервере при каждом входящем gRPC-запросе
func UnaryXRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, // Контекст запроса (живет пока обрабатывается запрос)
		req interface{}, // Входные данные (Proto-объект)
		info *grpc.UnaryServerInfo, // Информация о методе (какой метод вызвали)
		handler grpc.UnaryHandler, // Следующий обработчик в цепочке
	) (interface{}, error) {

		var requestID string

		md, ok := metadata.FromIncomingContext(ctx) //) для извлечения метаданных из контекста запроса. Метаданные — это аналог HTTP-заголовков.
		if ok {
			ids := md.Get("request-id")
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}

		//Если клиент передал свой request-id, код берет его. Если не передал — генерирует новый UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		return handler(ctx, req)

	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}