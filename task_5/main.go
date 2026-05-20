package main

import "fmt"

type User interface {
	GetUsername() string
	HasPermission(permission string) bool
	GetRole() string
}

type BasicUser struct {
	username string
	role     map[string]bool
}

type Moderator struct {
	BasicUser
}

type Admin struct {
	Moderator
}

func (name BasicUser) GetUsername() string {
	return name.username
}

func (name BasicUser) HasPermission(permission string) bool {
	return name.role[permission]
}

func (name BasicUser) GetRole() string {
	return "Basic"
}

func (m Moderator) GetRole() string {
	return "Moderator"
}

func (a Admin) GetRole() string {
	return "Admin"
}

func NewBasicUser(username string) BasicUser {
	return BasicUser{
		username: username,
		role: map[string]bool{
			"read": true,
		},
	}
}

func NewModerator(username string) Moderator {
	base := NewBasicUser(username)

	base.role["edit"] = true
	base.role["ban_user"] = true

	return Moderator{
		BasicUser: base,
	}
}

func NewAdmin(username string) Admin {
	mod := NewModerator(username)

	mod.role["delete"] = true
	mod.role["manage_roles"] = true

	return Admin{
		Moderator: mod,
	}
}

func main() {
	user := NewBasicUser("Анна")
	mod := NewModerator("Петя")
	admin := NewAdmin("Таня")

	users := []User{user, mod, admin}

	for _, u := range users {
		fmt.Println("Пользователь:", u.GetUsername(), "Роль:", u.GetRole())
		fmt.Println("  Может читать?", u.HasPermission("read"))
		fmt.Println("  Может редактировать?", u.HasPermission("edit"))
		fmt.Println("  Может удалять?", u.HasPermission("delete"))
	}
}