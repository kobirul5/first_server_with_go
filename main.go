package main

import (
	"context"
	"fmt"
	"net/http"

	"go-http-curd/db"
	"go-http-curd/handlers"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	db.ConnectDB()
	defer db.Db.Close(context.Background())

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/health", handlers.HandleHealth)
	mux.HandleFunc("/api/users", handlers.HandleUsers)
	mux.HandleFunc("POST /create-user", handlers.HandleCreateUser)
	mux.HandleFunc("POST /create-post", handlers.HandleCreatePost)
	mux.HandleFunc("PUT /update-user/{id}", handlers.HandleUpdateUser)
	mux.HandleFunc("PUT /update-post", handlers.HandleUpdatePost)
	mux.HandleFunc("DELETE /delete-user/{id}", handlers.HandleDeleteUser)

	// Dynamic route handling and parameters
	mux.HandleFunc("GET /user/{id}", handlers.HandleGetSingleUser)
	mux.HandleFunc("GET /api/users/search", handlers.HandleSearchUsers)
	mux.HandleFunc("POST /api/users/form", handlers.HandleCreateUserForm)
	mux.HandleFunc("GET /api/external-posts", handlers.HandleExternalData)

	fmt.Println("Server is running on port 5000...")
	err = http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
