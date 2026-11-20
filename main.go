package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "Rahim", Email: "rahim@example.com"},
	{ID: 2, Name: "Karim", Email: "karim@example.com"},
	{ID: 3, Name: "Kobirul", Email: "kobirul@example.com"},
}

func connectDB() {
	var err error

	
	connectStr := os.Getenv("DB_URL")

	db, err = pgx.Connect(context.Background(), connectStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}

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
	mux.HandleFunc("PUT /update-user", handleUpdateUser)
	mux.HandleFunc("PUT /update-user/{id}", handleUpdateUser)
	mux.HandleFunc("PUT /update-post", handleUpdatePost)
	mux.HandleFunc("PATCH /update-user", handlePatchUser)
	mux.HandleFunc("DELETE /delete-user/{id}", handleDeleteUser)

	// daynamic route handling and query parameters are not implemented in this example. You can extend the functionality as needed.
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

// get single user by id from database
func handleGetSingleUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, username, age, email FROM users WHERE id = $1`
	var user User
	err = db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Email,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

//---------------------------------------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// ServeMux matches all unmatched paths with "/", so handle 404 for non-root paths
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Welcome to the Go Server!")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health check passed!")
}

// Get all users from database
func handleUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, username, age, email FROM users ORDER BY id ASC`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userList = []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Email); err != nil {
			http.Error(w, "Failed to scan user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, u)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// Create users
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode incoming request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 2. Validate required fields (name and email)
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	// 3. Fallback ID assignment if not auto-generated
	if newUser.ID == 0 {
		maxID := 0
		for _, u := range users {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		newUser.ID = maxID + 1
	}

	// 4. Insert into database
	query := `
	insert into users (username, age, email)
	values ($1, $2, $3)
	returning id
	`
 
	err = db.QueryRow(context.Background(), query, newUser.Name, newUser.Age, newUser.Email).Scan(&newUser.ID)

	// Error handling
	if err != nil {
		http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return newly created user as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
	// json.NewEncoder(w).Encode(newUser)
}

// Create posts
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Post created successfully!")
}

// Replace/Update a user in database
func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ID from URL path or request body
	idParam := r.PathValue("id")
	var id int
	var err error

	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid user ID in path", http.StatusBadRequest)
			return
		}
	}

	// 2. Decode update data from request body
	var reqUser User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// If ID not in path, use ID from request body
	if id == 0 && reqUser.ID != 0 {
		id = reqUser.ID
	}

	if id == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// 3. Validation check
	if reqUser.Name == "" || reqUser.Email == "" {
		http.Error(w, "Username and Email are required", http.StatusBadRequest)
		return
	}

	// 4. Execute update query on database
	query := `
	UPDATE users
	SET username = $1, age = $2, email = $3
	WHERE id = $4
	RETURNING id, username, age, email
	`

	var updatedUser User
	err = db.QueryRow(context.Background(), query, reqUser.Name, reqUser.Age, reqUser.Email, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Age,
		&updatedUser.Email,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Send successful update response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// Replace/Update a post
func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Post updated successfully!")
}

// Partially update a user
func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User updated successfully!")
}

// Delete a user from database
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := db.Exec(context.Background(), query, id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted successfully!")
}

// 1. Data via URL Query Parameters (e.g. /api/users/search?username=Rahim&limit=5)
func handleSearchUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	limitParam := r.URL.Query().Get("limit")

	limit := 10
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	query := `SELECT id, username, age, email FROM users WHERE ($1 = '' OR username ILIKE '%' || $1 || '%') ORDER BY id ASC LIMIT $2`
	rows, err := db.Query(context.Background(), query, username, limit)
	if err != nil {
		http.Error(w, "Failed to search users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userList = []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Email); err != nil {
			http.Error(w, "Failed to scan user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

// 2. Data via HTML Form Data & Headers (instead of JSON body)
func handleCreateUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse URL-encoded / multipart form body
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Read data from request headers as well
	clientSource := r.Header.Get("X-Client-Source")
	if clientSource == "" {
		clientSource = "Web Form / HTML"
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")

	if username == "" || email == "" {
		http.Error(w, "Username and Email are required form fields", http.StatusBadRequest)
		return
	}

	age, _ := strconv.Atoi(ageStr)

	query := `
	INSERT INTO users (username, age, email)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var newID int
	err := db.QueryRow(context.Background(), query, username, age, email).Scan(&newID)
	if err != nil {
		http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":       "User created successfully via Form Data",
		"id":            newID,
		"username":      username,
		"email":         email,
		"age":           age,
		"client_source": clientSource,
	})
}

// 3. Data from an External 3rd-party API (No Database Used)
func handleExternalData(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?_limit=3")
	if err != nil {
		http.Error(w, "Failed to fetch external data: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var externalPosts []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&externalPosts); err != nil {
		http.Error(w, "Failed to parse external JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"source": "External 3rd-party API (No Database Used)",
		"posts":  externalPosts,
	})
}
