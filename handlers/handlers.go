package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-http-curd/db"
	"go-http-curd/model"

	"github.com/jackc/pgx/v5"
)

// RootHandler handles the root route "/"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// ServeMux matches all unmatched paths with "/", so handle 404 for non-root paths
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Welcome to the Go Server!")
}

// HandleHealth checks server health status
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health check passed!")
}

// HandleUsers fetches all users from database
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, username, age, email FROM users ORDER BY id ASC`

	rows, err := db.Db.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userList = []model.User{}
	for rows.Next() {
		var u model.User
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

// HandleGetSingleUser fetches a single user by ID from database
func HandleGetSingleUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, username, age, email FROM users WHERE id = $1`
	var user model.User
	err = db.Db.QueryRow(context.Background(), query, id).Scan(
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

// HandleCreateUser inserts a new user into the database
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode incoming request body into User struct
	var newUser model.User
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
		for _, u := range model.Users {
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

	err = db.Db.QueryRow(context.Background(), query, newUser.Name, newUser.Age, newUser.Email).Scan(&newUser.ID)

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
}

// HandleCreatePost creates a dummy post
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Post created successfully!")
}

// HandleUpdateUser updates an existing user in database
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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
	var reqUser model.User
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

	var updatedUser model.User
	err = db.Db.QueryRow(context.Background(), query, reqUser.Name, reqUser.Age, reqUser.Email, id).Scan(
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

// HandleUpdatePost updates a dummy post
func HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Post updated successfully!")
}

// HandleDeleteUser removes a user by ID from the database
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := db.Db.Exec(context.Background(), query, id)
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

// HandleSearchUsers searches users via URL Query Parameters
func HandleSearchUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	limitParam := r.URL.Query().Get("limit")

	limit := 10
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	query := `SELECT id, username, age, email FROM users WHERE ($1 = '' OR username ILIKE '%' || $1 || '%') ORDER BY id ASC LIMIT $2`
	rows, err := db.Db.Query(context.Background(), query, username, limit)
	if err != nil {
		http.Error(w, "Failed to search users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userList = []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Email); err != nil {
			http.Error(w, "Failed to scan user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

// HandleCreateUserForm creates a user via HTML Form Data & Headers
func HandleCreateUserForm(w http.ResponseWriter, r *http.Request) {
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
	err := db.Db.QueryRow(context.Background(), query, username, age, email).Scan(&newID)
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

// HandleExternalData fetches data from an external 3rd-party API
func HandleExternalData(w http.ResponseWriter, r *http.Request) {
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
