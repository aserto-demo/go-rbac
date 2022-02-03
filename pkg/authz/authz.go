package authz

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Authorizer interface {
	HasPermission(userID, action, asset string) bool
}

func Middleware(a Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			// This is where the password would normally be verified

			asset := mux.Vars(r)["asset"]
			action := actionFromMethod(r.Method)
			if !ok || !a.HasPermission(username, action, asset) {
				log.Printf("User '%s' is not allowed to '%s' resource '%s'", username, action, asset)
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
		return "gather"
	case "POST":
		return "consume"
	case "DELETE":
		return "destroy"
	default:
		return ""
	}
}
