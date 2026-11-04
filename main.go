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

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Course struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Instructor string  `json:"instructor"`
	Price      float64 `json:"price"`
}

type Order struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Item   string  `json:"item"`
	Total  float64 `json:"total"`
	Status string  `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/about", handleAbout)
	mux.HandleFunc("/contact", handleContact)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/products", handleProducts)
	mux.HandleFunc("/api/posts", handlePosts)
	mux.HandleFunc("/api/courses", handleCourses)
	mux.HandleFunc("/api/orders", handleOrders)

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

func handleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "About: This is a simple HTTP server written in Go.")
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Contact: Reach out to us at contact@example.com")
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

func handleProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 75000.0, Category: "Electronics"},
		{ID: 2, Name: "Mechanical Keyboard", Price: 3500.0, Category: "Accessories"},
		{ID: 3, Name: "Wireless Mouse", Price: 1500.0, Category: "Accessories"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	posts := []Post{
		{ID: 1, Title: "Introduction to Go", Author: "Kobirul"},
		{ID: 2, Title: "Building REST APIs in Go", Author: "Rahim"},
		{ID: 3, Title: "Concurrency in Go", Author: "Karim"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	courses := []Course{
		{ID: 1, Title: "Go Backend Development", Instructor: "Jhankar Mahbub", Price: 5000.0},
		{ID: 2, Title: "Full Stack Web Development", Instructor: "Programming Hero Team", Price: 6500.0},
		{ID: 3, Title: "Docker & Kubernetes Mastery", Instructor: "Kobirul Islam", Price: 4000.0},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	orders := []Order{
		{ID: 101, UserID: 1, Item: "Laptop", Total: 75000.0, Status: "Delivered"},
		{ID: 102, UserID: 2, Item: "Mechanical Keyboard", Total: 3500.0, Status: "Processing"},
		{ID: 103, UserID: 3, Item: "Go Backend Development", Total: 5000.0, Status: "Completed"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
