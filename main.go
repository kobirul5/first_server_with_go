package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "Rahim", Email: "rahim@example.com"},
	{ID: 2, Name: "Karim", Email: "karim@example.com"},
	{ID: 3, Name: "Kobirul", Email: "kobirul@example.com"},
}

func connectDB() {

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	var err error

	db, err = pgx.Connect(context.Background(), connectStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")

	// defer conn.Close(context.Background())

}

func main() {

	connectDB()

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("POST /create-user", handleCreateUser)
	mux.HandleFunc("POST /create-post", handleCreatePost)
	mux.HandleFunc("PUT /update-user", handleUpdateUser)
	mux.HandleFunc("PUT /update-post", handleUpdatePost)
	mux.HandleFunc("PATCH /update-user", handlePatchUser)
	mux.HandleFunc("DELETE /delete-user/{id}", handleDeleteUser)

	// daynamic route handling and query parameters are not implemented in this example. You can extend the functionality as needed.
	mux.HandleFunc("GET /user/{id}", handleGetSingleUser)

	fmt.Println("Server is running on port 5000...")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// get single user by id, dynamic route handling, and query parameters are not implemented in this example. You can extend the functionality as needed.

func handleGetSingleUser(w http.ResponseWriter, r *http.Request) {

	idParam := r.PathValue("id")
	fmt.Printf("the value or id %v and the type  of the id is %T", idParam, idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)

}

//---------------------------------------

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

	// ১. রিকোয়েস্ট বডি থেকে ডাটা রিসিভ করার জন্য User struct এর একটি ভ্যারিয়েবল তৈরি করা
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// ২. ভ্যালিডেশন চেক (নাম ও ইমেইল দেওয়া হয়েছে কি না)
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	// ৩. যদি ID না দেওয়া থাকে, তবে অটোমেটিক পরবর্তী ID নির্ধারণ করা
	if newUser.ID == 0 {
		maxID := 0
		for _, u := range users {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		newUser.ID = maxID + 1
	}

	// ৪. গ্লোবাল users স্লাইসে নতুন ইউজার যোগ করা
	users = append(users, newUser)

	// ৫. সফল রেসপন্স হিসেবে নতুন ইউজারকে রিটার্ন করা
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
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

// Replace/Update a post
func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Post updated successfully!")
}

// Partially update a user
func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User updated successfully!")
}

// Delete a user
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {

	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user ID")
		return
	}

	for idx, user := range users {
		if user.ID == id {
			users = slices.Delete(users, idx, idx+1)

			w.WriteHeader(http.StatusNoContent)
			fmt.Fprintln(w, "User deleted successfully!")
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "User not found")
}
