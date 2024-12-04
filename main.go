package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "tk_basdat"
	password = "123"
	dbname   = "tk_basdat"
)

var db *sql.DB

type LoginRequestBody struct {
	NoHP string `json:"NoHP"`
	Pwd  string `json:"Pwd"`
}

type LoginResponseBody struct {
	Status  bool   `json:"status"`
	Role    int    `json:"role"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type RegisterRequestBody struct {
	Role              int       `json:"role"`
	Nama              string    `json:"name"`
	JenisKelamin      string    `json:"sex"`
	NoHP              string    `json:"number"`
	Pwd               string    `json:"password"`
	TglLahir          time.Time `json:"date"`
	Alamat            string    `json:"address"`
	NamaBank          string    `json:"bank"`
	NomorRekening     string    `json:"noRek"`
	NPWP              string    `json:"npwp"`
	LinkFoto          string    `json:"link"`
	Rating            float64   `json:"rating"`
	JmlPsnananSelesai int       `json:"amount"`
}

type RegisterResponseBody struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func main() {
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	db = conn
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Connected to the PostgreSQL database")

	// tambah endpoint disini
	http.HandleFunc("/login", corsMiddleware(checkLogin))
	http.HandleFunc("/register", corsMiddleware(register))

	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userId string
	err = db.QueryRow(`INSERT INTO "user" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id`, uuid.New(), body.Nama, body.JenisKelamin, body.NoHP, body.Pwd, body.TglLahir, body.Alamat, 0.0).Scan(&userId)
	if err == sql.ErrNoRows {
		response := &RegisterResponseBody{
			Status:  false,
			Message: "Invalid Credential on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &RegisterResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	if body.Role == 0 {
		err = db.QueryRow(`INSERT INTO PELANGGAN VALUES ($1, $2) RETURNING Id`, userId, "Basic").Scan(&userId)
		if err == sql.ErrNoRows {
			response := &RegisterResponseBody{
				Status:  false,
				Message: "Invalid Credential on pelanggan",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &RegisterResponseBody{
				Status:  false,
				Message: err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		err = db.QueryRow(`INSERT INTO PEKERJA VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING Id`,
			userId,
			body.NamaBank,
			body.NomorRekening,
			body.NPWP,
			body.LinkFoto,
			body.Rating,
			body.JmlPsnananSelesai).Scan(&userId)

		if err == sql.ErrNoRows {
			response := &RegisterResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &RegisterResponseBody{
				Status:  false,
				Message: err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := &RegisterResponseBody{
		Status:  true,
		Message: "User berhasil dibuat",
	}

	json.NewEncoder(w).Encode(response)
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID string
	var name string
	var role int

	err = db.QueryRow(`SELECT Id, Nama FROM "user" WHERE NoHP = $1 AND Pwd = $2`, body.NoHP, body.Pwd).Scan(&userID, &name)
	if err == sql.ErrNoRows {
		response := &LoginResponseBody{
			Status:  false,
			UserId:  userID,
			Name:    name,
			Role:    role,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &LoginResponseBody{
			Status:  false,
			UserId:  userID,
			Name:    name,
			Role:    role,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	db.QueryRow(`SELECT 1 FROM PELANGGAN WHERE Id = $1`, userID).Scan(&role)
	response := &LoginResponseBody{
		Status:  true,
		UserId:  userID,
		Name:    name,
		Role:    role,
		Message: "Success",
	}

	json.NewEncoder(w).Encode(response)
}

// func addUser(w http.ResponseWriter, r *http.Request) {
//     var user User
//     err := json.NewDecoder(r.Body).Decode(&user)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     _, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusCreated)
//     fmt.Fprintf(w, "User added successfully")
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
//     var user User
//     err := json.NewDecoder(r.Body).Decode(&user)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     _, err = db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", user.Name, user.Email, user.ID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     fmt.Fprintf(w, "User updated successfully")
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
//     id := r.URL.Query().Get("id")
//     if id == "" {
//         http.Error(w, "ID parameter is required", http.StatusBadRequest)
//         return
//     }

//     _, err := db.Exec("DELETE FROM users WHERE id=$1", id)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     fmt.Fprintf(w, "User deleted successfully")
// }
