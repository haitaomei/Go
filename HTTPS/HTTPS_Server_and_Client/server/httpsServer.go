package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	/*------------ Restful style server ----------------*/
	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")
	helloRouter := router.PathPrefix("/helloAPI").Subrouter()
	helloRouter.HandleFunc("/{name}", helloHandler).Methods("GET")

	http.Handle("/", router)

	/* the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate. */
	http.ListenAndServeTLS(":60443", "../tls/cert.pem", "../tls/key.pem", nil) //https
	// http.ListenAndServe(":60443", nil)	//standard http
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
