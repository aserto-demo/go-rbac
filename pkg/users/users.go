package users

import (
	"github.com/aserto-demo/go-rbac/pkg/file"
)

type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

type Users map[string]User

func Load() (Users, error) {
	var userList []User

	if err := file.LoadJson("../users.json", &userList); err != nil {
		return nil, err
	}

	users := Users{}
	for _, user := range userList {
		users[user.ID] = user
	}

	return users, nil
}