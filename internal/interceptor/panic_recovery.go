package interceptor

import (
	"context"
	"log"

	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryPanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
	// recover() - ловит панику. Если паники не было, вернет nil
		if r := recover(); r != nil {
			// debug.Stack() - получаем стек вызовов для отладки
			log.Printf("PANIC RECOVERED: %v\nStack: %s", r, debug.Stack())

			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
	}
}

