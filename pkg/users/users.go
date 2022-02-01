package users

import (
	"encoding/json"

	"github.com/aserto-demo/go-rbac/pkg/file"
)

type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

type Users map[string]User

func Load() (Users, error) {
	jsonBytes, err := file.ReadBytes("../users.json")
	if err != nil {
		return nil, err
	}

	var userList []User
	if err := json.Unmarshal(jsonBytes, &userList); err != nil {
		return nil, err
	}

	users := Users{}
	for _, user := range userList {
		users[user.ID] = user
	}

	return users, nil
}
