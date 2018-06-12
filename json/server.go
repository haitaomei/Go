package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

/*
 * Create this using pointers, to test if a the client acutaly sent the id field.
 * Because if id is 0, we can not tell if the client sent zero, or omited it.
 */

type UserForRead struct {
	Name *string `json:"name"`
	ID   *int    `json:"id"`
}

func rootHandler(httpResp http.ResponseWriter, httpReq *http.Request) {
	var u UserForRead
	if httpReq.Body == nil {
		http.Error(httpResp, "Please send a request body", 400)
		return
	}
	defer httpReq.Body.Close()

	err := json.NewDecoder(httpReq.Body).Decode(&u)
	if err != nil {
		http.Error(httpResp, err.Error(), 400)
		return
	}

	fmt.Println(u)

	httpResp.Header().Add("Content-Type", "application/json")
	httpResp.WriteHeader(200)
	if u.ID == nil {
		json.NewEncoder(httpResp).Encode("ID is omited in the request")
	} else {
		json.NewEncoder(httpResp).Encode("ID is set")
	}
}

func main() {
	router := mux.NewRouter()
	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.HandleFunc("/", rootHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":60443", nil)
}
