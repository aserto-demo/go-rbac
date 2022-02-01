package main

import (
	"encoding/json"

	"github.com/aserto-demo/go-rbac/pkg/file"
)

type Resources []string

type Actions map[string]Resources

type Roles map[string]Actions

func LoadRoles() (Roles, error) {
	jsonBytes, err := file.ReadBytes("./roles.json")
	if err != nil {
		return nil, err
	}

	var roles Roles
	if err := json.Unmarshal(jsonBytes, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}
