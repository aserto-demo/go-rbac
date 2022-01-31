package main

import (
	"fmt"
	"log"
	"net/http"

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

	authz := authorizer{users: users, enforcer: enforcer}

	router := mux.NewRouter()
	router.Use(authz.Middleware)
	router.HandleFunc("/api/{asset}", HandleRequest).Methods("GET", "POST", "DELETE")

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

func (a *authorizer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()
		// This is where the password would normally be verified

		asset := mux.Vars(r)["asset"]
		action := actionFromMethod(r.Method)
		if !ok || !a.hasPermission(username, action, asset) {
			log.Printf("User '%s' is not allowed to '%s' resource '%s'", username, action, asset)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *authorizer) hasPermission(userID, action, asset string) bool {
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

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("Got permission"))
}
