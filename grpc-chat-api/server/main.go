package main

import (
	"context"
	pb "grpc-chat-api/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
}

func (c *ChatServer) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	log.Println("Сервер получил запрос на создание ")
	return &pb.CreateChatResponse{Id: 1}, nil
}

func (c *ChatServer) DeleteChat(ctx context.Context, req *pb.DeleteChatRequest) (*emptypb.Empty, error){
	log.Println("Сервер выполняет delete-запрос ID ", req.Id)
	return &emptypb.Empty{},nil
}

func (c *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error){
	log.Println("Сервер выполняет отправку сообщения на сервер ", req.Text, req.From, req.Timestamp)
	return &emptypb.Empty{},nil
}

func main(){
	lis, err := net.Listen("tcp", "localhost:3001")
	if err != nil{
		log.Fatalln("Не удалось подключиться к порту ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(grpcServer, &ChatServer{})

	log.Println("Сервер запущен и слушает порт 3001..")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}