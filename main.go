package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()	

	mux.HandleFunc("/", rootHandler)
	// mux.HandleFunc("/api", apiCallHandler)

	fmt.Println("Server is running on port 5000...")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

