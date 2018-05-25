package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Https Server Running")
}

func main() {
	/*------------ basic server ----------------*/
	// http.HandleFunc("/", handler)

	/*------------ Restful style server ----------------*/
	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")
	helloRouter := router.PathPrefix("/helloAPI").Subrouter()
	helloRouter.HandleFunc("/{name}", helloHandler).Methods("GET")

	http.Handle("/", router)

	http.ListenAndServeTLS(":60443", "server.crt", "server.key", nil)
}

func rootHandler(httpResp http.ResponseWriter, httpReq *http.Request) {

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode("Welcome to the root directory ...")
}

func helloHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	vars := mux.Vars(httpReq)
	username := vars["name"]
	var responseText = "Hi " + username + ", how are you?"

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	json.NewEncoder(httpResp).Encode(responseText)
}
