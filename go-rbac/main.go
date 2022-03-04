package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/gorilla/mux"
	"github.com/mikespook/gorbac"
)

func LoadJson(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func SaveJson(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

func main() {
	// map[RoleId]PermissionIds
	var jsonRoles map[string][]string

	// Load roles information
	if err := LoadJson("roles.json", &jsonRoles); err != nil {
		log.Fatal(err)
	}

	rbac := gorbac.New()
	permissions := make(gorbac.Permissions)

	// Build roles and add them to goRBAC instance
	for rid, pids := range jsonRoles {
		role := gorbac.NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = gorbac.NewStdPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac.Add(role)
	}

	users, err := users.Load()
	if err != nil {
		log.Fatal("Failed to load users:", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/{asset}", server.Handler).Methods("GET", "POST", "DELETE")
	router.Use(
		authz.Middleware(&authorizer{users: users, rbac: rbac, permissions: permissions}),
	)

	server.Start(router)

}
type authorizer struct {
	users    users.Users
	rbac *gorbac.RBAC
	permissions gorbac.Permissions
}


func (a *authorizer) HasPermission(userID, action, asset string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, role := range user.Roles {
		permission := action + "-" + asset
		if a.rbac.IsGranted(role, a.permissions[permission], nil) {
			return true
		}
	}

	return false
}