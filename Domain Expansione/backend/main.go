package main

import (
	"backend/handlers"
	"backend/redisclient"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	redisclient.InitRedis()
	// DB connection
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	// http.HandleFunc("/login", handlers.ShowLoginPage)
	http.Handle("/login/submit", handlers.LoginHandler(db))
	http.Handle("/register", handlers.RegisterHandler(db))
	http.Handle("/order/confirm", handlers.OrderSubmitHandler(db))
	http.Handle("/payment/confirm", handlers.OrderConfirmHandler(db))

	// Start server
	fmt.Println("Server running at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
