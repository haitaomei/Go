package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Https Server Running")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":60443", "server.crt", "server.key", nil)
}
