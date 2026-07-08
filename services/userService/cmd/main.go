package main

import (
    "log"
    "net"

    "google.golang.org/grpc"

    userv1 "grpc-exchange/gen/user"
    "github.com/grpc-exchange/services/userService/internal/handler"
    "github.com/grpc-exchange/shared/interceptor"
    "github.com/grpc-exchange/shared/logger"
)

func main() {
  if err := logger.InitLogger(); err != nil {
    log.Fatalf("Failed to init logger: %v", err)
  }

  lis, err := net.Listen("tcp", ":50053")
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  s := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
      interceptor.XRequestIDUnaryInterceptor(),
      interceptor.LoggerUnaryInterceptor(),
      interceptor.PanicRecoveryUnaryInterceptor(),
    ),
  )

  userHandler := handler.NewUserHandler()
  userv1.RegisterUserServiceServer(s, userHandler)

  log.Println("UserService listening on :50053")
  if err := s.Serve(lis); err != nil {
    log.Fatalf("Failed to serve: %v", err)
  }
}