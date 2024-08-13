package server

import (
	"fmt"
	"log"
	"net/http"
)

const Port = 8000

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`"Access granted"`))
}

func Start(handler http.Handler) {
	addr := fmt.Sprintf("0.0.0.0:%d", Port)
	fmt.Println("Staring server on", addr)

	srv := http.Server{
		Handler: handler,
		Addr:    addr,
	}
	log.Fatal(srv.ListenAndServe())
}
