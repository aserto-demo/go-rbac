package main

import (
	"log"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/gorilla/mux"
	"github.com/samber/lo"
)

type authorizer struct {
	users users.Users
	roles Roles
}

func (a *authorizer) HasPermission(userID, action, resource string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, roleName := range user.Roles {
		role := a.roles[roleName]
		if role == nil {
			log.Printf("User '%s' has unknown role '%s'", userID, roleName)
			continue
		}

		if allowed, ok := role[action]; ok {
			if lo.Contains(allowed, resource) {
				return true
			}
		}
	}

	return false
}

func main() {
	users, err := users.Load()
	if err != nil {
		log.Fatal("Failed to load users:", err)
	}

	roles, err := LoadRoles()
	if err != nil {
		log.Fatal("Failed to load roles:", err)
	}

	router := mux.NewRouter()
	router.Use(
		authz.Middleware(&authorizer{users: users, roles: roles}),
	)

	router.HandleFunc("/api/{resource}", server.Handler).Methods("GET", "PUT", "DELETE")

	server.Start(router)
}
