package main

import (
	"fmt"
	"sync"
	"time"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://victoriametrics.com/blog/go-sync-once/

//Не всегда понимаю, какие поля должны быть у структуры
type Comment struct{
	id int
	text string
	userId int
}
type User struct{
	id int
	name string
}
type Sesion struct{
	sessionId string
	userId int
}
type File struct{
	id int
	url string
	commentId int
}

func loadComment(comments *[]Comment, wg *sync.WaitGroup){
	defer wg.Done()
	
	fmt.Println("Загрузка комментариев")
	time.Sleep(200*time.Millisecond)
	
	*comments = []Comment{
		{id: 1, text: "One", userId: 1},
	}
	
	fmt.Println("Комментарии загружены")

}

func (s *Sesion) loadSession(session *Sesion, wg *sync.WaitGroup){
	fmt.Println("Загрузка session")
	defer wg.Done()
	time.Sleep(100*time.Millisecond)
	
	*session = Sesion{
		sessionId: "abc",
		userId: 1,
	}
	fmt.Println("Session загружены")
}

func loadUser(userId int, users *map[int]User, wg *sync.WaitGroup){
	defer wg.Done()
	time.Sleep(100 * time.Millisecond)

	(*users)[userId] = User{
		id: userId,
		name: string(userId),
	}

	fmt.Println("Users загружены")
}

func main(){
	wg := sync.WaitGroup{}
	wgCom := sync.WaitGroup{}

	var comments []Comment
	var session Sesion
	// тк будем искать пользователя по id, чтобы находить мгновенно
	users := make(map[int]User)
	
	sission := &Sesion{}

	wg.Add(1)
	go loadComment(&comments, &wg)
	
	wg.Add(1)
	go sission.loadSession(&session, &wg)

	wgCom.Wait()
	fmt.Println("Comment загружены")

	fmt.Println("загрузка пользователей")
	//проходимся по всем комментариям и собираем userId
	unicUser := make(map[int]bool)
	for _, comm := range comments{
		unicUser[comm.id] = true
	}

	// Запускаем загрузку каждого пользователя
	for userID := range unicUser{
		wg.Add(1)
		go loadUser(userID, &users, &wg)
	}
	wg.Wait()

	fmt.Println("Sessions", session.sessionId)
	fmt.Println("Comments", comments)
	fmt.Println("Users", users)
}