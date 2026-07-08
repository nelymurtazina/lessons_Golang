package interceptor

import (
    "context"
    "runtime/debug"

    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "github.com/grpc-exchange/shared/logger"
)

func PanicRecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
        defer func() {
            if r := recover(); r != nil {
                logger.Log.Error("Panic recovered",
                    zap.Any("panic", r),
                    zap.String("stack", string(debug.Stack())),
                    zap.String("method", info.FullMethod),
                )
                err = status.Errorf(codes.Internal, "internal server error: %v", r)
            }
        }()
        return handler(ctx, req)
    }
}