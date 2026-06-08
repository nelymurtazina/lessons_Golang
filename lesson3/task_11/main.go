package main

import (
	"fmt"
	"sync"
	"time"
	"strconv"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://victoriametrics.com/blog/go-sync-once/

type Comment struct{
	ID int
	Text string
	UserId int
}
type User struct{
	ID int
	Name string
}
type Session struct{
	SessionId string
	UserId int
}

func loadComments(comments *[]Comment, wg *sync.WaitGroup){
	defer wg.Done()
	
	time.Sleep(200*time.Millisecond)

	*comments = []Comment{
		{ID: 1,Text: "Первый", UserId: 101},
	}

	fmt.Println("Комментарии загружены")
}

func loadSession(session *Session, wg *sync.WaitGroup){
	defer wg.Done()

	time.Sleep(100*time.Millisecond)

	*session = Session{
		SessionId: "session-123",
		UserId: 101,
	}

	fmt.Println("Сессия загружены")
}

func loadUser(userId int, users *map[int]User, mu *sync.Mutex, wg *sync.WaitGroup){
	defer wg.Done()

	time.Sleep(100*time.Millisecond)

	mu.Lock()
	(*users)[userId] = User{
		ID: userId,
		Name: "User_" + strconv.Itoa(userId),
	}
	mu.Unlock()

	fmt.Println("Загружен пользователь")
}

func main(){
	var wg sync.WaitGroup
	var mu sync.Mutex

	var comments []Comment
	var session Session
	users := make(map[int]User)

	wg.Add(2)
	go loadComments(&comments, &wg)
	go loadSession(&session, &wg)
	wg.Wait()

	uniqueUserIDs := make(map[int]bool)
	for _, comcomments := range comments{
		uniqueUserIDs[comcomments.UserId] = true
	}

	wg.Add(len(uniqueUserIDs))
	for userID := range uniqueUserIDs{
		go loadUser(userID, &users,&mu,&wg)
	}
	wg.Wait()

	fmt.Printf("Сессия: %s\n", session.SessionId)
	fmt.Printf("Пользователей: %d\n", len(users))
	for _, u := range users {
		fmt.Printf("   - %d: %s\n", u.ID, u.Name)
	}
	fmt.Printf("Комментариев: %d\n", len(comments))
}