package handler

import (
  "context"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  userv1 "grpc-exchange/gen/user"
  "github.com/grpc-exchange/services/userService/internal/domain"
)

type UserHandler struct {
  userv1.UnimplementedUserServiceServer
  users map[string]*domain.User
}

func NewUserHandler() *UserHandler {
  // Тестовые пользователи
  users := map[string]*domain.User{
    "user-123": {
      ID:       "user-123",
      Username: "alice",
      Email:    "alice@example.com",
      Roles:    []string{"user", "trader"},
      Active:   true,
      },
      "admin-456": {
        ID:       "admin-456",
        Username: "bob",
        Email:    "bob@example.com",
        Roles:    []string{"admin"},
        Active:   true,
      },
  }

  return &UserHandler{
    users: users,
  }
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error) {
  user, ok := h.users[req.UserId]
  if !ok {
    return nil, status.Errorf(codes.NotFound, "user %s not found", req.UserId)
  }

  return &userv1.User{
    UserId:   user.ID,
    Username: user.Username,
    Email:    user.Email,
    Roles:    user.Roles,
    Active:   user.Active,
  }, nil
}

func (h *UserHandler) CheckUserRoles(ctx context.Context, req *userv1.CheckUserRolesRequest) (*userv1.CheckUserRolesResponse, error) {
  user, ok := h.users[req.UserId]
  if !ok {
    return nil, status.Errorf(codes.NotFound, "user %s not found", req.UserId)
  }

  // Проверяем, есть ли у пользователя хотя бы одна из требуемых ролей
  hasAccess := false
  for _, requiredRole := range req.RequiredRoles {
    for _, userRole := range user.Roles {
      if userRole == requiredRole {
        hasAccess = true
        break
      }
    }
    if hasAccess {
      break
    }
  }

  return &userv1.CheckUserRolesResponse{
    HasAccess:  hasAccess,
    UserRoles:  user.Roles,
  }, nil
}