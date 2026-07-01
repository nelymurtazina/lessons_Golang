package interceptor

import (
    "context"
    "time"

    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/status"
)

func UnaryLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		requestID := GetRequestID(ctx)
		code := status.Code(err)

		logger.Info("gRPC call completed",
			zap.String("method", info.FullMethod), // например "/order.OrderService/CreateOrder"
			zap.String("request_id", requestID),
			zap.Duration("duration", duration),
			zap.String("status", code.String()),
			zap.Error(err),
		)
		return resp, err
	}
}