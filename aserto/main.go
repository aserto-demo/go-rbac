package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-dev/aserto-go/authorizer/grpc"
	"github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/middleware"
	"github.com/aserto-dev/aserto-go/middleware/http/std"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func AsertoAuthorizer(addr, tenantID, apiKey, policyID string) (*std.Middleware, error) {
	ctx := context.Background()
	authClient, err := grpc.New(
		ctx,
		client.WithAddr(addr),
		client.WithTenantID(tenantID),
		client.WithAPIKeyAuth(apiKey),
	)
	if err != nil {
		return nil, err
	}

	mw := std.New(
		authClient,
		middleware.Policy{
			ID:       policyID,
			Decision: "allowed",
		},
	)
	mw.Identity.Mapper(func(r *http.Request, identity middleware.Identity) {
		if username, _, ok := r.BasicAuth(); ok {
			identity.Subject().ID(username)
		}
	})
	mw.WithPolicyFromURL("gorbac")
	return mw, nil
}

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	authorizerAddr := os.Getenv("AUTHORIZER_ADDRESS")
	if authorizerAddr == "" {
		authorizerAddr = "authorizer.prod.aserto.com:8443"
	}
	apiKey := os.Getenv("AUTHORIZER_API_KEY")
	policyID := os.Getenv("POLICY_ID")
	tenantID := os.Getenv("TENANT_ID")

	authorizer, err := AsertoAuthorizer(authorizerAddr, tenantID, apiKey, policyID)
	if err != nil {
		log.Fatal("Failed to create authorizer:", err)
	}

	log.Print(os.Getenv("AUTHORIZER_API_KEY"))

	router := mux.NewRouter()
	router.HandleFunc("/api/{asset}", server.Handler).Methods("GET", "POST", "DELETE")
	router.Use(authorizer.Handler)

	server.Start(router)
}