package interceptor

import (
    "context"

    "github.com/google/uuid"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

type ctxKeyRequestID struct{}

func XRequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        var requestID string

        if md, ok := metadata.FromIncomingContext(ctx); ok {
            ids := md.Get("x-request-id")
            if len(ids) > 0 {
                requestID = ids[0]
            }
        }

        if requestID == "" {
            requestID = uuid.New().String()
        }

        ctx = context.WithValue(ctx, ctxKeyRequestID{}, requestID)
        ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)

        return handler(ctx, req)
    }
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
        return id
    }
    return ""
}


