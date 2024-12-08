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
func main() {
	// Connect to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	// Define routes
	http.HandleFunc("/homepage", getHomepage)
	http.HandleFunc("/subkategori", getSubkategori)
	http.HandleFunc("/pesan", createPesanan)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Get Homepage (all categories and subcategories)
func getHomepage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT k.id, k.nama, s.id, s.nama
		FROM kategori_jasa k
		LEFT JOIN subkategori_jasa s ON k.id = s.id_kategori`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var kategoriID, subkategoriID int
		var kategoriNama, subkategoriNama string
		if err := rows.Scan(&kategoriID, &kategoriNama, &subkategoriID, &subkategoriNama); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"kategori_id":   kategoriID,
			"kategori_nama": kategoriNama,
			"subkategori_id": subkategoriID,
			"subkategori_nama": subkategoriNama,
		})
	}
	json.NewEncoder(w).Encode(data)
}

// Get Subkategori and Sessions
func getSubkategori(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	rows, err := db.Query(`
		SELECT s.id, s.nama, s.deskripsi, sesi.id, sesi.nama_sesi, sesi.harga
		FROM subkategori_jasa s
		LEFT JOIN sesi_layanan sesi ON s.id = sesi.id_subkategori
		WHERE s.id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var subID, sesiID int
		var subNama, subDeskripsi, sesiNama string
		var harga float64
		if err := rows.Scan(&subID, &subNama, &subDeskripsi, &sesiID, &sesiNama, &harga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"subkategori_id": subID,
			"subkategori_nama": subNama,
			"subkategori_deskripsi": subDeskripsi,
			"sesi_id": sesiID,
			"sesi_nama": sesiNama,
			"harga": harga,
		})
	}
	json.NewEncoder(w).Encode(data)
}

// Create Order
func createPesanan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		UserID           string  `json:"user_id"`
		SesiID           int     `json:"sesi_id"`
		Tanggal          string  `json:"tanggal"`
		Diskon           float64 `json:"diskon"`
		MetodePembayaran string  `json:"metode_pembayaran"`
		Total            float64 `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
		INSERT INTO pesanan (id_user, id_sesi, tanggal, diskon, metode_pembayaran, total, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'Menunggu Pembayaran')`,
		body.UserID, body.SesiID, body.Tanggal, body.Diskon, body.MetodePembayaran, body.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Order created"))
}