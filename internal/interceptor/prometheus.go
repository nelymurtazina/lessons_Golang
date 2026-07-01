package interceptor

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func UnaryPrometheusInterceptor() grpc.UnaryServerInterceptor {
    return grpc_prometheus.UnaryServerInterceptor
}