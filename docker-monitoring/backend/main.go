package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "mydb"
)

var db *sql.DB

// PingResult структура для приема данных от pinger
type PingResult struct {
	IP        string    `json:"ip"`
	Latency   string    `json:"latency"`
	Timestamp time.Time `json:"timestamp"`
}

// savePing сохраняет данные пинга в БД
func savePing(result PingResult) error {
	_, err := db.Exec("INSERT INTO pings (ip, latency, timestamp) VALUES ($1, $2, $3)", result.IP, result.Latency, result.Timestamp)
	return err
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var result PingResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Ошибка обработки JSON", http.StatusBadRequest)
		return
	}

	if err := savePing(result); err != nil {
		http.Error(w, "Ошибка сохранения данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные успешно сохранены"))
}

func main() {
	// Подключение к базе
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	}

	fmt.Println("✅ Подключение к базе успешно!")

	// Создаём таблицу, если её нет
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pings (
		id SERIAL PRIMARY KEY,
		ip TEXT NOT NULL,
		latency TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL
	)`)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}

	// HTTP-сервер
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Backend работает!"))
	})
	http.HandleFunc("/ping", handlePing) // Добавляем обработчик

	fmt.Println("🚀 Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
