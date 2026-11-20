# Go HTTP Server (CRUD API)

A simple RESTful HTTP server built using Go's standard library (`net/http`), Go 1.22+ `http.ServeMux` routing, and PostgreSQL (`pgx/v5`).

---

## 📁 Project Structure

```text
go_server/
├── explore/
│   └── explore.go      # Go slice manipulation practice file
├── .gitignore          # Git ignore file
├── GEMINI.md           # Project guidelines & conventions
├── go.mod              # Go module file
├── go.sum              # Go dependencies checksums
├── main.go             # Server setup, database connection, routes, and handlers
└── README.md           # Project documentation
```

---

## ⚙️ How It Works

### 1. Server Setup (`main.go`)
- Uses `http.NewServeMux()` as the request router.
- Connects to PostgreSQL using `github.com/jackc/pgx/v5`.
- Listens and serves on port `:5000` (`http.ListenAndServe(":5000", mux)`).

### 2. Routing Mechanism (Go 1.22+ Routing)
Starting with Go 1.22, HTTP methods and dynamic path parameters are supported directly in patterns:
- **Method Matching:** e.g., `POST /create-user`, `GET /user/{id}`, `DELETE /delete-user/{id}`.
- **Path Parameters:** `r.PathValue("id")` extracts wildcard values from the URL path.

### 3. Database Integration (PostgreSQL)
User data is stored and managed in a PostgreSQL database using `pgx/v5`:
- **Create:** `INSERT INTO users ... RETURNING id`
- **Read (All):** `SELECT id, username, age, email FROM users ORDER BY id ASC`
- **Read (Single):** `SELECT id, username, age, email FROM users WHERE id = $1`
- **Update:** `UPDATE users SET ... WHERE id = $4 RETURNING ...`
- **Delete:** `DELETE FROM users WHERE id = $1`

---

## 🚀 API Endpoints

| Method | Endpoint | Description | Status Code |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Root handler (welcome message) | `200 OK` |
| `GET` | `/health` | Health check endpoint | `200 OK` |
| `GET` | `/api/users` | Returns list of all users from database (JSON) | `200 OK` |
| `GET` | `/user/{id}` | Returns a single user by ID from database | `200 OK` / `404 Not Found` |
| `POST` | `/create-user` | Creates a new user in database | `201 Created` |
| `POST` | `/create-post` | Creates a new post | `200 OK` |
| `PUT` | `/update-user/{id}` | Updates an existing user in database | `200 OK` / `404 Not Found` |
| `PUT` | `/update-post` | Updates a post | `200 OK` |
| `PATCH` | `/update-user` | Partially updates a user | `200 OK` |
| `DELETE` | `/delete-user/{id}` | Deletes a user by ID from database | `200 OK` / `404 Not Found` |
| `GET` | `/api/users/search` | Search users via URL Query Parameters (`?username=...&limit=...`) | `200 OK` |
| `POST` | `/api/users/form` | Create user via Form Data & Headers (instead of JSON body) | `201 Created` |
| `GET` | `/api/external-posts` | Fetches data from an external 3rd-party API (No DB used) | `200 OK` |

---

## 🛠️ How to Run the Server

Run the server from the project root directory:

```bash
go run .
```

When started successfully, the console will display:
```text
Database connected successfully
Server is running on port 5000...
```

---

## 🧪 Testing with cURL

#### 1. Get All Users
```bash
curl http://localhost:5000/api/users
```

#### 2. Get Single User by ID (e.g. ID: 1)
```bash
curl http://localhost:5000/user/1
```

#### 3. Create a User (POST)
```bash
curl -X POST http://localhost:5000/create-user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Abdullah",
    "email": "abdullah@example.com",
    "age": 28
  }'
```

#### 4. Update a User (PUT)
```bash
curl -X PUT http://localhost:5000/update-user/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Abdullah Al Mamun",
    "email": "abdullah.mamun@example.com",
    "age": 29
  }'
```

#### 5. Delete a User (DELETE)
```bash
curl -X DELETE http://localhost:5000/delete-user/1
```

#### 6. Health Check
```bash
curl http://localhost:5000/health
```

---

## 🔍 `explore/explore.go`
Demonstrates slice manipulation techniques for removing elements:
```go
// Removes the element at index `idx`:
users = append(users[:idx], users[idx+1:]...)
```
