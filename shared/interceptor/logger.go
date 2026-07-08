package interceptor

import (
    "context"
    "time"

    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/status"

    "github.com/grpc-exchange/shared/logger"
)

func LoggerUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        start := time.Now()
        requestID := GetRequestID(ctx)

        resp, err := handler(ctx, req)
        duration := time.Since(start)

        fields := []zap.Field{
            zap.String("method", info.FullMethod),
            zap.Duration("duration", duration),
            zap.String("request_id", requestID),
        }

        if err != nil {
            st, _ := status.FromError(err)
            fields = append(fields, zap.String("error", st.Message()), zap.Int("code", int(st.Code())))
            logger.Log.Error("Request failed", fields...)
        } else {
            logger.Log.Info("Request completed", fields...)
        }

        return resp, err
    }
}