package authz

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Authorizer interface {
	HasPermission(userID, action, resource string) bool
}

func Middleware(a Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			// This is where the password would normally be verified

			resource := mux.Vars(r)["resource"]
			action := actionFromMethod(r.Method)
			if !ok || !a.HasPermission(username, action, resource) {
				log.Printf("User '%s' is not allowed to '%s' resource '%s'", username, action, resource)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func actionFromMethod(httpMethod string) string {
	switch httpMethod {
	case "GET":
		return "can_read"
	case "PUT":
		return "can_write"
	case "DELETE":
		return "can_delete"
	default:
		return ""
	}
}
