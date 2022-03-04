package main

import (
	"log"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/gorilla/mux"
)

type authorizer struct {
	users users.Users
	roles Roles
}

func (a *authorizer) HasPermission(userID, action, asset string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, roleName := range user.Roles {
		if role, ok := a.roles[roleName]; ok {
			resources, ok := role[action]
			if ok {
				for _, resource := range resources {
					if resource == asset {
						return true
					}
				}
			}
		} else {
			log.Printf("User '%s' has unknown role '%s'", userID, roleName)
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
	router.HandleFunc("/api/{asset}", server.Handler).Methods("GET", "POST", "DELETE")
	router.Use(
		authz.Middleware(&authorizer{users: users, roles: roles}),
	)

	server.Start(router)
}