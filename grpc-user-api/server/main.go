package main

import (
	"context"
	"fmt"
	pb "grpc-user-api/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct{
  // Встраиваем сгенерированный интерфейс. обязательно
  pb.UnimplementedUserServiceServer
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
  fmt.Println("Сервер получил запрос на создание пользователя", req.Name, req.Email,req.Role.String())
  
  return &pb.CreateResponse{Id: 123}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error){
  fmt.Println("Сервер выполняет get-запрос")
  return &pb.GetResponse{
    Id: req.Id,
    Name: "Петя Петров",
    Email: "petr@gamil.com",
    Role: pb.Role_ROLE_USER,
    CreatedAt: timestamppb.Now(),
    UpdateAt: timestamppb.Now(),
  }, nil
}

func (s *UserServer) UpdateUser (ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error){
  fmt.Println("Сервер выполняет update-запрос")

  if req.Name != nil {
    fmt.Println("Новое имя: ", req.Name.GetValue())
  }

  if req.Email != nil {
    fmt.Println("Новое email: ", req.Email.GetValue())
  }

  return &emptypb.Empty{},nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error){
  fmt.Println("Сервер выполняет delete-запрос, на пользователя с ID ", req.Id)

  return &emptypb.Empty{},nil
}


func main(){
  lis, err := net.Listen("tcp", "localhost:50051")
  if err != nil {
    log.Fatalln("Не удалось подключиться к порту ", err)
  }

  // Создаем gRPC сервер
	grpcServer := grpc.NewServer()

  // Регистрируем UserServer в gRPC-сервере.
	// Мы говорим gRPC: "Все запросы к UserService отправляй "
	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	fmt.Println("Сервер запущен и слушает порт 50051...")
	
	// Запускаем сервер (он будет работать бесконечно)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}