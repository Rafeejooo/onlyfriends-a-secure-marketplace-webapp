package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func OrderSubmitHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		name := r.FormValue("name")
		// talentID := r.FormValue("talent_id")
		duration := r.FormValue("package")
		date := r.FormValue("date")

		if name == "" || duration == "" || date == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO orders (name, package, date)
			VALUES ($1, $2, $3)
			RETURNING id
		`

		var orderID int
		err := db.QueryRow(query, name, duration, date).Scan(&orderID)
		if err != nil {
			// 🚨 Log the actual DB error!
			log.Println("INSERT error:", err) // <-- Add this
			http.Error(w, "Failed to submit order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
	})
}

func OrderConfirmHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		orderID := r.FormValue("order_id")

		query := `
			UPDATE orders SET status = 'confirmed' WHERE id = $1
		`

		_, err := db.Exec(query, orderID)
		if err != nil {
			http.Error(w, "Failed to confirm order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
