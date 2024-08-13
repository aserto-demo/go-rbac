package main

import (
	"log"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
)

type authorizer struct {
	users    users.Users
	enforcer *casbin.Enforcer
}

func (a *authorizer) HasPermission(userID, action, resource string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, role := range user.Roles {
		if a.enforcer.Enforce(role, resource, action) {
			return true
		}
	}

	return false
}

func main() {
	enforcer, err := casbin.NewEnforcerSafe("./rbac_model.conf", "./rbac_policy.csv")
	if err != nil {
		log.Fatal("Failed to create enforcer:", err)
	}

	users, err := users.Load()
	if err != nil {
		log.Fatal("Failed to load users:", err)
	}

	router := mux.NewRouter()
	router.Use(
		authz.Middleware(&authorizer{users: users, enforcer: enforcer}),
	)

	router.HandleFunc("/api/{resource}", server.Handler).Methods("GET", "PUT", "DELETE")

	server.Start(router)
}
