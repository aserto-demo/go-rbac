package users

import (
	"encoding/json"
	"io"
	"os"
)

type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

type Users map[string]User

func Load() (Users, error) {
	jsonBytes, err := readFileBytes("../users.json")
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

func readFileBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
