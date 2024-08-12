package main

import (
	"log"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/file"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-demo/go-rbac/pkg/users"
	"github.com/gorilla/mux"
	"github.com/mikespook/gorbac"
)

type authorizer struct {
	users       users.Users
	rbac        *gorbac.RBAC
	permissions gorbac.Permissions
}

func (a *authorizer) HasPermission(userID, action, resource string) bool {
	user, ok := a.users[userID]
	if !ok {
		// Unknown userID
		log.Print("Unknown user:", userID)
		return false
	}

	for _, role := range user.Roles {
		permission := action + "-" + resource
		if a.rbac.IsGranted(role, a.permissions[permission], nil) {
			return true
		}
	}

	return false
}

func main() {
	// map[RoleId]PermissionIds
	var roles map[string][]string

	// Load roles information
	if err := file.LoadJson("roles.json", &roles); err != nil {
		log.Fatal(err)
	}

	rbac := gorbac.New()
	permissions := make(gorbac.Permissions)

	// Build roles and add them to goRBAC instance
	for rid, pids := range roles {
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
	router.HandleFunc("/api/{resource}", server.Handler).Methods("GET", "PUT", "DELETE")
	router.Use(
		authz.Middleware(&authorizer{users: users, rbac: rbac, permissions: permissions}),
	)

	server.Start(router)
}
