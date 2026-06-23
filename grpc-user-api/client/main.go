package main

import (
	"context"
	"log"
	"time"

	pb "grpc-user-api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main(){
	//подключение к серверу
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}

	defer conn.Close()
	//Создаем "клиентскую заглушку"
	client := pb.NewUserServiceClient(conn)

	// Контекст с таймаутом (чтобы клиент не вис вечно, если сервер упал)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Клиент: отправляю запрос на создание")
	createResp, err := client.CreateUser(ctx, &pb.CreateRequest{
		Name: "Nely",
		Email: "nmur@gmail.ru",
		Password: "1234567", //а так норм хранить пароль? я просто не знаю, как по-другому 
		PasswordConfirm: "1234567",
		Role: pb.Role_ROLE_ADMIN,
	})
	if err != nil{
		log.Fatal("Ошибка создания ", err)
	}
	log.Println("Новый ID: ", createResp.Id)

	log.Println("Клиент: Отправляю запрос на получение")
	getResp, err := client.GetUser(ctx, &pb.GetRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Ошибка Get: %v", err)  // ← Исправлено!
	}
	log.Printf("Клиент получил ответ: Name=%s, Email=%s\n", getResp.Name, getResp.Email)


	log.Println(" Клиент: Отправляю запрос на обновление")
	// Для StringValue нужно использовать специальную обертку wrapperspb.String()
	//wrapperspb.String(): в Update у нас StringValue. Это сделано для того, чтобы можно было передать nil (не обновлять поле). В Go для этого нужно использовать специальные обертки из wrapperspb.
	_, err = client.UpdateUser(ctx, &pb.UpdateRequest{
		Id:    createResp.Id,
		Name:  wrapperspb.String("Nely new"),
		Email: nil, // email не меняем
	})
	if err != nil {
		log.Fatalf("Ошибка Update: %v", err)
	}
	log.Println("Клиент: Update успешно выполнен (пустой ответ)!")



	log.Println("Клиент: Отправляю запрос на удаление")
	_, err = client.DeleteUser(ctx, &pb.DeleteRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Ошибка Delete: %v", err)
	}
	log.Println("Клиент: Delete успешно выполнен (пустой ответ)!")
}