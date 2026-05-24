package main

import "fmt"

var User map[string]int

func init(){
	User = make(map[string]int)
}

func GetAge(name string) int {
	fmt.Println(User[name])
	return User[name]
}

func DeletePerson(name string){
	if _, ok := User[name]; !ok{
		fmt.Println("Такое человка нет в списке ", name)
	} 
	delete(User, name)
}

func PrintAll(){
	for value := range User{
		fmt.Println("Имя: ", value,", Возраст: ", User[value])
	}
}

func main() {
	User["Petr"] = 30
	User["Alise"] = 21

	GetAge("Petr")

	fmt.Println("ДО: ",User)
	DeletePerson("Nely")
	DeletePerson("Alise")
	fmt.Println("После: ",User)
	User["Nikola"] = 30
	User["Kate"] = 21
	PrintAll()
}