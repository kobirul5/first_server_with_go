package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	connectDB()
	defer db.Close(context.Background())

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("POST /create-user", handleCreateUser)
	mux.HandleFunc("POST /create-post", handleCreatePost)
	mux.HandleFunc("PUT /update-user/{id}", handleUpdateUser)
	mux.HandleFunc("PUT /update-post", handleUpdatePost)
	mux.HandleFunc("DELETE /delete-user/{id}", handleDeleteUser)

	// Dynamic route handling and parameters
	mux.HandleFunc("GET /user/{id}", handleGetSingleUser)
	mux.HandleFunc("GET /api/users/search", handleSearchUsers)
	mux.HandleFunc("POST /api/users/form", handleCreateUserForm)
	mux.HandleFunc("GET /api/external-posts", handleExternalData)

	fmt.Println("Server is running on port 5000...")
	err = http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
