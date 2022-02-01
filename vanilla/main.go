package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/gorilla/mux"
)

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
	router.HandleFunc("/api/{asset}", handleRequest).Methods("GET", "POST", "DELETE")
	router.Use(
		authz.Middleware(&authorizer{users: users, roles: roles}),
	)

	fmt.Println("Staring server on 0.0.0.0:8080")

	srv := http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

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

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("Got permission"))
}

func actionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "view"
	case "POST":
		return "edit"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}
