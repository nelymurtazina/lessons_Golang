package main

import (
	"fmt"
	"sync"
	"time"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://victoriametrics.com/blog/go-sync-once/

type Comment struct{
	Id string
	AuthorID string
	Text string
}

type User struct{
	Id string
	Name string
}

type Sessia struct{
	Id string
	IsValid   bool
}

type File struct{
	CommentID string
	FileURL   string
}

// loadComments имитирует загрузку комментариев из БД
func loadComments() []Comment{
	time.Sleep(500*time.Millisecond)
	return []Comment{
		{Id: "1", AuthorID: "user1", Text: "Первый комментарий"},
	}
}

func loadSesia() Sessia{
	time.Sleep(300*time.Millisecond)
	return Sessia{
		Id: "session1",
		IsValid: true,
	}
}

func loadUser(userId string) User{
	time.Sleep(100 * time.Millisecond)
	return User{Id: userId, Name: "Пользователь " + userId}
}

func loadFile(comments []Comment, sessionId string){
	fmt.Println("Загрузка сессии: ", sessionId)
	time.Sleep(200*time.Millisecond)
	fmt.Println("Вложения загружены")
}

func main(){
	//Независимые Комментарии и сессии
	comments := []Comment{}
	session := Sessia{}
	users := make(map[string]User)
	wg := sync.WaitGroup{}
	once := sync.Once{}

	wg.Add(2)

	go func(){
		defer wg.Done()
		fmt.Println("Загрузка")
		comments = loadComments()
		fmt.Println("Загрузка комментариев: ", len(comments))
	}()

	go func(){
		defer wg.Done()
		fmt.Println("Загрузка")
		session = loadSesia()
		if session.IsValid{
			fmt.Println("Загрузка сессии: ", session.Id)
		} else{
			fmt.Println("Сессия не найдена")
		}
	}()

	wg.Wait()

	wg.Add(1)
  go func() {
    defer wg.Done()
		// Собираем уникальные ID авторов
    userIDs := make(map[string]bool)

    for _, c := range comments {
      userIDs[c.AuthorID] = true
    }

		for userID := range userIDs {
  	user := loadUser(userID)
    users[userID] = user
    fmt.Println("Загружен пользователь: ", user.Name, user.Id)
  }
	}()

	wg.Wait()

	if session.IsValid && session.Id != "" {
		once.Do(func() {
			fmt.Println("Загрузка вложений ")
			loadFile(comments, session.Id)
		})
	} else {
		fmt.Println("Нет вложений")
	}

	fmt.Println("Комментарии:")
	for _, c := range comments {
    user := users[c.AuthorID]
    fmt.Printf(c.Id, user.Name, c.Text)
  }
}