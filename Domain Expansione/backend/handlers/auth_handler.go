package handlers

import (
	"backend/auth"
	"backend/jwt"
	"backend/redisclient"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("pswd")

		// Authenticate user from DB
		id, username, err := auth.AuthenticateUser(db, email, password)
		if err != nil || id == 0 {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}

		// Cache the username in Redis after successful login
		err = redisclient.Rdb.Set(redisclient.Ctx, email, username, 24*time.Hour).Err()
		if err != nil {
			// Log this, but don’t fail the login flow if Redis is unavailable
			fmt.Println("Redis SET error:", err)
		}

		// Generate JWT with username
		token, err := jwt.GenerateJWT(email, username)
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		//  Set cookie with JWT
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}
}

func RegisterHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("pswd")
		username := r.FormValue("txt")
		phone := r.FormValue("broj")

		err := auth.RegisterUser(db, email, password, username, phone)
		if err != nil {
			http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}
