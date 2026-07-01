package service

import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var (
    // Ошибки для Market
    ErrMarketNotFound = status.Error(codes.NotFound, "market not found")
    ErrMarketNotActive = status.Error(codes.FailedPrecondition, "market is not active")
    
    // Ошибки для Order
    ErrOrderNotFound = status.Error(codes.NotFound, "order not found")
    ErrInvalidArgument = status.Error(codes.InvalidArgument, "invalid argument")
    ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
)