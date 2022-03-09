package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-dev/aserto-go/authorizer/grpc"
	"github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/middleware"
	"github.com/aserto-dev/aserto-go/middleware/http/std"
)

var ErrMissingVar = errors.New("missing environment variable(s)")

type asertoEnv struct {
	addr     string
	apiKey   string
	policyID string
	tenantID string
}

func (env *asertoEnv) Validate() error {
	missing := []string{}
	if env.addr == "" {
		missing = append(missing, "AUTHORIZER_ADDRESS")
	}
	if env.apiKey == "" {
		missing = append(missing, "AUTHORIZER_API_KEY")
	}
	if env.policyID == "" {
		missing = append(missing, "POLICY_ID")
	}
	if env.tenantID == "" {
		missing = append(missing, "TENANT_ID")
	}

	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingVar, missing)
	}

	return nil
}

func AsertoAuthorizer(env *asertoEnv) (*std.Middleware, error) {
	ctx := context.Background()
	authClient, err := grpc.New(
		ctx,
		client.WithAddr(env.addr),
		client.WithTenantID(env.tenantID),
		client.WithAPIKeyAuth(env.apiKey),
	)
	if err != nil {
		return nil, err
	}

	mw := std.New(
		authClient,
		middleware.Policy{
			ID:       env.policyID,
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

func loadEnv() (*asertoEnv, error) {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: failed to load .env file.")
	}

	authorizerAddr := os.Getenv("AUTHORIZER_ADDRESS")
	if authorizerAddr == "" {
		authorizerAddr = "authorizer.prod.aserto.com:8443"
	}

	env := &asertoEnv{
		addr:     authorizerAddr,
		apiKey:   os.Getenv("AUTHORIZER_API_KEY"),
		policyID: os.Getenv("POLICY_ID"),
		tenantID: os.Getenv("TENANT_ID"),
	}

	if err := env.Validate(); err != nil {
		return nil, err
	}

	return env, nil
}

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatalf("Environment error: %s", err)
	}

	authorizer, err := AsertoAuthorizer(env)
	if err != nil {
		log.Fatal("Failed to create authorizer:", err)
	}

	log.Print(os.Getenv("AUTHORIZER_API_KEY"))

	router := mux.NewRouter()
	router.HandleFunc("/api/{asset}", server.Handler).Methods("GET", "POST", "DELETE")
	router.Use(authorizer.Handler)

	server.Start(router)
}
