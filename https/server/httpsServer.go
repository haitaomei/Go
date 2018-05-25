package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Https Server Running")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":60443", "server.crt", "server.key", nil)

	r := mux.NewRouter()
	inuseRouter := r.PathPrefix("/api/v1/project").Subrouter()
}
