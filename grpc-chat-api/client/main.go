package main

import (
	"context"
	"log"
	"time"

	pb "grpc-chat-api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main(){
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}

	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Client: create chat")
	createResp, err := client.CreateChat(ctx, &pb.CreateChatRequest{
		Usernames: []string{"Alice", "Kate", "Nick"},
	})
	if err != nil {
		log.Fatal("Ошибка создания", err)
	}
	log.Println("Новый ID chat: ", createResp.Id)

	log.Println("Client: send message")
	_, err = client.SendMessage(ctx, &pb.SendMessageRequest{
		From:"Alice",
		Text:"Привет!",
		Timestamp: timestamppb.Now(),  
	})
	if err != nil {
		log.Fatalf("Ошибка SendMessage: %v", err)
	}
	log.Println("Сообщение отправлено")


	log.Println("Client: delete chat")
	//_ означает "игнорируем первый возврат"
	_, err = client.DeleteChat(ctx, &pb.DeleteChatRequest{Id: createResp.Id})
	if err != nil{
		log.Fatalf("Ошибка Delete: %v", err)
	}
	log.Println("delete успешно выполнен (пустой ответ)")
}