package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("POST /create-user", handleCreateUser)
	mux.HandleFunc("POST /create-post", handleCreatePost)
	mux.HandleFunc("PUT /update-user", handleUpdateUser)
	mux.HandleFunc("PATCH /update-user", handlePatchUser)
	mux.HandleFunc("DELETE /delete-user", handleDeleteUser)

	fmt.Println("Server is running on port 5000...")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// ServeMux এ "/" সব unmatched পাথে ম্যাচ করে, তাই 404 হ্যান্ডেল করা হলো
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Welcome to the Go Server!")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health check passed!")
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "Rahim", Email: "rahim@example.com"},
		{ID: 2, Name: "Karim", Email: "karim@example.com"},
		{ID: 3, Name: "Kobirul", Email: "kobirul@example.com"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// Create users
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "User created successfully!")
}

// Create posts
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Post created successfully!")
}

// Replace a user
func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User replaced successfully!")
}

// Partially update a user
func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User updated successfully!")
}

// Delete a user
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User deleted successfully!")
}
