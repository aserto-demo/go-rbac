package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
)

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
	router.HandleFunc("/api/{asset}", handleRequest).Methods("GET", "POST", "DELETE")
	router.Use(
		authz.Middleware(&authorizer{users: users, enforcer: enforcer}),
	)

	fmt.Println("Staring server on 0.0.0.0:8080")

	srv := http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

type authorizer struct {
	users    users.Users
	enforcer *casbin.Enforcer
}

func (a *authorizer) HasPermission(userID, action, asset string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, role := range user.Roles {
		if a.enforcer.Enforce(role, asset, action) {
			return true
		}
	}

	return false
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("Got permission"))
}
