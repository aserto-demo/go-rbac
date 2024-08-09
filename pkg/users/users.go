package users

import (
	"github.com/aserto-demo/go-rbac/pkg/file"
	"github.com/samber/lo"
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

	users := lo.Associate(userList, func(u User) (string, User) {
		return u.ID, u
	})

	return users, nil
}
