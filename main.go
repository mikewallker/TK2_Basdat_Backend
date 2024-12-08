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

type GetUserRequestBody struct {
	User string `json:"user"`
	Role int    `json:"role"`
}

type GetUserResponseBody struct {
	Status            bool      `json:"status"`
	Message           string    `json:"message"`
	User              string    `json:"userid"`
	Role              int       `json:"role"`
	Nama              string    `json:"name"`
	JenisKelamin      string    `json:"sex"`
	NoHP              string    `json:"number"`
	Pwd               string    `json:"password"`
	TglLahir          time.Time `json:"date"`
	Alamat            string    `json:"address"`
	SaldoMyPay        float64   `json:"saldo"`
	Level             string    `json:"level"`
	NamaBank          string    `json:"bank"`
	NomorRekening     string    `json:"noRek"`
	NPWP              string    `json:"npwp"`
	LinkFoto          string    `json:"link"`
	Rating            float64   `json:"rating"`
	JmlPsnananSelesai int       `json:"amount"`
}

type UpdateUserRequestBody struct {
	User          string    `json:"user"`
	Role          int       `json:"role"`
	Nama          string    `json:"name"`
	JenisKelamin  string    `json:"sex"`
	NoHP          string    `json:"number"`
	TglLahir      time.Time `json:"date"`
	Alamat        string    `json:"address"`
	NamaBank      string    `json:"bank"`
	NomorRekening string    `json:"noRek"`
	NPWP          string    `json:"npwp"`
	LinkFoto      string    `json:"link"`
}

type UpdateUserResponseBody struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

//BAGIAN MERAH
type MyPayHistory struct {
	ID        string  `json:"id"`
	Tgl       string  `json:"date"`
	Nominal   float64 `json:"nominal"`
	Kategori  string  `json:"category"`
}

type MyPayHistoryResponse struct {
	UserID   string         `json:"user_id"`
	History  []MyPayHistory `json:"history"`
}

type MyPayTransactionTransfer struct {
	UserID     string  `json:"user_id"`
	KategoriID string  `json:"kategori_id"`
	Nominal    float64 `json:"nominal"`
	ToUserID   string  `json:"to_user_id,omitempty"` // For transfers
}

type MyPayTransactionTopUp struct {
	UserID    string  `json:"userId"`
	Nominal   float64 `json:"nominal"`
	KategoriID int    `json:"kategoriId"`
}

type MyPayKategori struct{
	NamaKategori string `json:"namaKategori"`
}

type GetPesananJasaRequestBody struct {
	User string `json:"user"`
}

type GetPesananJasaResponseBody struct {
	Status   bool            `json:"status"`
	Message  string          `json:"message"`
	Pesanan  []PesananJasa   `json:"pesanan"`
}

type PesananJasa struct {
	NamaJasa  string  `json:"nama_jasa"`
	TotalBiaya float64 `json:"total_biaya"`
}

type MyPayTransactionPay struct {
	UserID     string  `json:"user_id"`
	KategoriID string  `json:"kategori_id"`
	Nominal    float64 `json:"nominal"`
	ToUserID   string  `json:"to_user_id,omitempty"` // For transfers
	BankName   string  `json:"bank_name,omitempty"` // For withdrawals
	AccountNo  string  `json:"account_no,omitempty"`
}

type MyPayTransactionWithdrawal struct {
	UserID     string  `json:"user_id"`
	KategoriID string  `json:"kategori_id"`
	Nominal    float64 `json:"nominal"`
	BankName   string  `json:"bank_name,omitempty"` // For withdrawals
	AccountNo  string  `json:"account_no,omitempty"`
}


type AvailableJobsRequest struct {
	WorkerID string `json:"worker_id"`
}

type Job struct {
	JobID        string  `json:"job_id"`
	ServiceName  string  `json:"service_name"`
	ScheduledAt  string  `json:"scheduled_at"`
	TotalCost    float64 `json:"total_cost"`
}

type AvailableJobsResponse struct {
	Jobs []Job  `json:"jobs"`
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
	http.HandleFunc("/getUser", corsMiddleware(getUser))
	http.HandleFunc("/updateUser", corsMiddleware(updateUser))
	//ENDPOINT MERAH
	http.HandleFunc("/mypay/balance", corsMiddleware(getMyPayBalance))
	http.HandleFunc("/mypay/history", corsMiddleware(getMyPayHistory))
	http.HandleFunc("/mypay/topup", corsMiddleware(handleTopUp))
	http.HandleFunc("mypay/get-category-id", corsMiddleware(GetCategoryIdByName))
	http.HandleFunc("mypay/getPesananJasa", corsMiddleware(getPesananJasa))
	http.HandleFunc("mypay/getStatusIdByName", corsMiddleware(GetStatusIdByName))
	http.HandleFunc("mypay/processPayment", corsMiddleware(ProcessPayment))
	http.HandleFunc("/mypay/transaction", corsMiddleware(handleMyPayTransaction))
	http.HandleFunc("/mypay/transaction", corsMiddleware(handleMyPayTransaction))
	http.HandleFunc("/mypay/transaction", corsMiddleware(handleMyPayTransaction))
	http.HandleFunc("/jobs/available", corsMiddleware(getAvailableJobs))

	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
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

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body UpdateUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var oldValue UpdateUserRequestBody
	err = db.QueryRow(`SELECT Nama, JenisKelamin, NoHP, TglLahir, Alamat FROM "user" WHERE Id = $1`, body.User).
		Scan(
			&oldValue.Nama,
			&oldValue.JenisKelamin,
			&oldValue.NoHP,
			&oldValue.TglLahir,
			&oldValue.Alamat)

	if err == sql.ErrNoRows {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: "Invalid Credential on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	var current_user_id string
	err = db.QueryRow(`UPDATE "user" SET Nama = $1, JenisKelamin = $2, TglLahir = $3, Alamat = $4 WHERE Id = $5 Returning Id`,
		oldValue.Nama,
		oldValue.JenisKelamin,
		oldValue.TglLahir,
		oldValue.Alamat,
		body.User).Scan(&current_user_id)

	if body.NoHP != oldValue.NoHP {
		err = db.QueryRow(`UPDATE "user" SET NoHP = $1 WHERE Id = $2 Returning Id`,
			body.NoHP,
			body.User).Scan(&current_user_id)
	}

	if err == sql.ErrNoRows {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: "Invalid update on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: err.Error() + " User",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	if body.Role == 1 {
		err = db.QueryRow(`SELECT NPWP, LinkFoto, NamaBank, NomorRekening FROM PEKERJA WHERE Id = $1`, body.User).
			Scan(
				&oldValue.NPWP,
				&oldValue.LinkFoto,
				&oldValue.NamaBank,
				&oldValue.NomorRekening)
		if err == sql.ErrNoRows {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: err.Error() + " Update",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		err = db.QueryRow(`UPDATE PEKERJA SET 
        NPWP = $1, 
        LinkFoto = $2 
        WHERE Id = $3 Returning Id`,
			oldValue.NPWP,
			oldValue.LinkFoto,
			body.User).Scan(&current_user_id)
		if err == sql.ErrNoRows {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: err.Error() + " Pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		if body.NomorRekening != oldValue.NomorRekening && body.NamaBank != oldValue.NamaBank {
			fmt.Println(body.NomorRekening + " " + body.NamaBank + " " + oldValue.NomorRekening + " " + oldValue.NamaBank)
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NamaBank = $1, 
			NomorRekening = $2 
			WHERE Id = $3 Returning Id`,
				body.NamaBank,
				body.NomorRekening,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Update",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		} else if body.NamaBank != oldValue.NamaBank {
			fmt.Println("Masuk-[0]")
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NamaBank = $1
			WHERE Id = $2 Returning Id`,
				body.NamaBank,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Bank Name",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		} else if body.NomorRekening != oldValue.NomorRekening {
			fmt.Println("Masuk")
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NomorRekening = $1
			WHERE Id = $2 Returning Id`,
				body.NomorRekening,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Rekening",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		}

	}

	response := &RegisterResponseBody{
		Status:  true,
		Message: fmt.Sprint("User dengan id $s berhasil di update", current_user_id),
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

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body GetUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response GetUserResponseBody
	err = db.QueryRow(`SELECT Nama, JenisKelamin, NoHP, Pwd, TglLahir, Alamat, SaldoMyPay FROM "user" WHERE Id = $1`, body.User).Scan(
		&response.Nama,
		&response.JenisKelamin,
		&response.NoHP,
		&response.Pwd,
		&response.TglLahir,
		&response.Alamat,
		&response.SaldoMyPay)

	if err == sql.ErrNoRows {
		response := &GetUserResponseBody{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &GetUserResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Message = "Berhasil mendapatkan data"

	if body.Role == 0 {
		db.QueryRow(`SELECT Level FROM PELANGGAN WHERE Id = $1`, body.User).Scan(&response.Level)
		json.NewEncoder(w).Encode(response)
	} else {
		db.QueryRow(`SELECT NamaBank, NomorRekening, NPWP, LinkFoto, Rating, JmlPsnananSelesai FROM PEKERJA WHERE Id = $1`, body.User).Scan(
			&response.NamaBank,
			&response.NomorRekening,
			&response.NPWP,
			&response.LinkFoto,
			&response.Rating,
			&response.JmlPsnananSelesai)
		json.NewEncoder(w).Encode(response)
	}
}

//BAGIAN MERAH
// MyPay model to track user's balance
type MyPay struct {
    ID       int     `json:"id"`
    UserID   int     `json:"user_id"`
    Balance  float64 `json:"balance"`
}

// Service Order model
type ServiceOrder struct {
    ID           int    `json:"id"`
    UserID       int    `json:"user_id"`
    WorkerID     int    `json:"worker_id"`
    ServiceName  string `json:"service_name"`
    Status       string `json:"status"`  // Example: "Looking for Worker", "Worker Assigned", etc.
    CreatedAt    string `json:"created_at"`
    ScheduledAt  string `json:"scheduled_at"`
}


// Get MyPay Balance for User
func getMyPayBalance(w http.ResponseWriter, r *http.Request) {
    var request GetUserRequestBody

    // Decode the JSON body to get user ID
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if request.User == "" {
        http.Error(w, "Missing user ID in request body", http.StatusBadRequest)
        return
    }

    var balance float64
	var noHP string 

    // Use the provided user ID to get the balance
    err := db.QueryRow("SELECT SaldoMyPay, NoHP FROM \"user\" WHERE Id = $1", request.User).Scan(&balance, &noHP)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
        return
    }

    // Set the response headers and write the balance in JSON format
    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "user_id": request.User,
        "balance": balance,
		"no_hp":   noHP,
    }
    json.NewEncoder(w).Encode(response)
}

// Get History Transaction for User
func getMyPayHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Change to POST if you're sending data in the request body
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body into the GetUserRequestBody struct
	var requestBody GetUserRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if user ID is present in the request body
	if requestBody.User == "" {
		http.Error(w, "Missing user ID in request body", http.StatusBadRequest)
		return
	}

	// Query the transaction history using the user ID from the request body
	rows, err := db.Query(`
		SELECT t.Id, t.Tgl, t.Nominal, k.Nama 
		FROM TR_MYPAY t
		JOIN KATEGORI_TR_MYPAY k ON t.KategoriId = k.Id
		WHERE t.UserId = $1
		ORDER BY t.Tgl DESC`, requestBody.User)
	if err != nil {
		http.Error(w, "Failed to fetch transaction history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []MyPayHistory
	for rows.Next() {
		var transaction MyPayHistory
		err := rows.Scan(&transaction.ID, &transaction.Tgl, &transaction.Nominal, &transaction.Kategori)
		if err != nil {
			http.Error(w, "Failed to parse transaction history", http.StatusInternalServerError)
			return
		}
		history = append(history, transaction)
	}

	// Respond with the transaction history in the required format
	response := MyPayHistoryResponse{
		UserID:  requestBody.User,
		History: history,
	}
	json.NewEncoder(w).Encode(response)
}

// Handle Top-Up
func handleTopUp(w http.ResponseWriter, r *http.Request) {
	var request MyPayKategori
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var transaction MyPayTransactionTopUp
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE \"user\" SET SaldoMyPay = SaldoMyPay + $1 WHERE Id = $2", transaction.Nominal, transaction.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to top-up", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`INSERT INTO TR_MYPAY (Id, UserId, Tgl, Nominal, KategoriId) VALUES (uuid_generate_v4(), $1, CURRENT_DATE, $2, $3)`,
		transaction.UserID, transaction.Nominal, transaction.KategoriID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to record top-up transaction", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Top-up successful"})
}

// GetCategoryIdByName fetches the category UUID based on the category name
func GetCategoryIdByName(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse the request body into the MyPayKategori struct
    var requestBody MyPayKategori
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&requestBody); err != nil {
        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
        return
    }

    // Validate the category name
    if requestBody.NamaKategori == "" {
        http.Error(w, "Category name is required", http.StatusBadRequest)
        return
    }

    // Query to get the category UUID based on the category name
    var kategoriId uuid.UUID
    err := db.QueryRow(`
        SELECT Id
        FROM KATEGORI_TR_MYPAY
        WHERE Nama = $1`, requestBody.NamaKategori).Scan(&kategoriId)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Category not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to fetch category UUID", http.StatusInternalServerError)
        }
        return
    }

    // Respond with the category UUID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "kategoriId": kategoriId.String(),
    })
}

func getPesananJasa(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body GetPesananJasaRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response GetPesananJasaResponseBody

	rows, err := db.Query(`
		SELECT 
			k.NamaKategori, 
			tpj.TotalBiaya 
		FROM 
			TR_PEMESANAN_JASA tpj
		JOIN 
			KATEGORI_JASA k ON k.Id = tpj.IdKategoriJasa
		JOIN 
			TR_PEMESANAN_STATUS tps ON tpj.Id = tps.IdTrPemesanan
		WHERE 
			tps.Keterangan = 'Menunggu Pembayaran'
			AND tpj.IdPelanggan = $1`, body.User)
	if err != nil {
		response.Status = false
		response.Message = "Error executing query: " + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var pesanan []PesananJasa
	for rows.Next() {
		var pesananItem PesananJasa
		err := rows.Scan(&pesananItem.NamaJasa, &pesananItem.TotalBiaya)
		if err != nil {
			response.Status = false
			response.Message = "Error scanning row: " + err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		pesanan = append(pesanan, pesananItem)
	}

	// Pastikan response selalu mengembalikan pesanan, meskipun kosong
	response.Status = true
	if len(pesanan) == 0 {
		response.Message = "No orders found"
	} else {
		response.Message = "Successfully retrieved orders"
	}
	response.Pesanan = pesanan

	json.NewEncoder(w).Encode(response)
}

func GetStatusIdByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parsing request body
	var requestBody struct {
		StatusName string `json:"statusName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if requestBody.StatusName == "" {
		http.Error(w, "Status name is required", http.StatusBadRequest)
		return
	}

	// Fetch status ID from STATUS_PESANAN
	var statusId uuid.UUID
	err := db.QueryRow(`SELECT Id FROM STATUS_PESANAN WHERE Status = $1`, requestBody.StatusName).Scan(&statusId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Status not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch status ID", http.StatusInternalServerError)
		}
		return
	}

	// Respond with status ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"statusId": statusId.String(),
	})
}


func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parsing request body
	var requestBody struct {
		UserId    string `json:"userId"`    // UUID format
		ServiceId string `json:"serviceId"` // UUID format
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate UUID format for UserId and ServiceId
	if _, err := uuid.Parse(requestBody.UserId); err != nil {
		http.Error(w, "Invalid UserId format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(requestBody.ServiceId); err != nil {
		http.Error(w, "Invalid ServiceId format", http.StatusBadRequest)
		return
	}

	// Fetch service price
	var servicePrice float64
	err := db.QueryRow(`SELECT TotalBiaya FROM TR_PEMESANAN_JASA WHERE Id = $1`, requestBody.ServiceId).Scan(&servicePrice)
	if err != nil {
		http.Error(w, "Failed to fetch service price", http.StatusInternalServerError)
		return
	}

	// Fetch user balance
	var userBalance float64
	err = db.QueryRow(`SELECT SaldoMyPay FROM "user" WHERE Id = $1`, requestBody.UserId).Scan(&userBalance)
	if err != nil {
		http.Error(w, "Failed to fetch user balance", http.StatusInternalServerError)
		return
	}

	// Validate balance
	if userBalance < servicePrice {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Saldo tidak mencukupi untuk melakukan pembayaran.",
		})
		return
	}

	// Update user balance
	_, err = db.Exec(`UPDATE "user" SET SaldoMyPay = SaldoMyPay - $1 WHERE Id = $2`, servicePrice, requestBody.UserId)
	if err != nil {
		http.Error(w, "Failed to update user balance", http.StatusInternalServerError)
		return
	}

	// Fetch new status ID
	var newStatusId uuid.UUID
	err = db.QueryRow(`SELECT Id FROM STATUS_PESANAN WHERE Status = 'Mencari Pekerja Terdekat'`).Scan(&newStatusId)
	if err != nil {
		http.Error(w, "Failed to fetch new status ID", http.StatusInternalServerError)
		return
	}

	// Update service status
	_, err = db.Exec(`
		UPDATE TR_PEMESANAN_JASA 
		SET IdKategoriJasa = $1, IdDiskon = NULL, IdMetodeBayar = NULL 
		WHERE Id = $2`, newStatusId, requestBody.ServiceId)
	if err != nil {
		http.Error(w, "Failed to update service status", http.StatusInternalServerError)
		return
	}

	// Insert into TR_MYPAY
	var categoryId uuid.UUID
	err = db.QueryRow(`SELECT Id FROM KATEGORI_TR_MYPAY WHERE Nama = 'membayar transaksi jasa'`).Scan(&categoryId)
	if err != nil {
		http.Error(w, "Failed to fetch category ID", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO TR_MYPAY (Id, UserId, Tgl, Nominal, KategoriId) VALUES (uuid_generate_v4(), $1, CURRENT_DATE, $2, $3)`,
		requestBody.UserId, servicePrice, categoryId)
	if err != nil {
		http.Error(w, "Failed to insert transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Payment successful",
	})
}

