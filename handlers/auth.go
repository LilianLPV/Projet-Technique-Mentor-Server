package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/config"
	"server/middleware"
	"server/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur lors du hash", http.StatusInternalServerError)
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		string(hashedPassword),
	)
	var storedHash string
	config.DB.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedHash)
	err2 := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	log.Println("Test bcrypt immédiat:", err2)
	if err != nil {
		log.Println("Erreur register:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("User créé avec succès")
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string

	err = config.DB.QueryRow(
		"SELECT id_user, username, email, password FROM users WHERE email = ?",
		credentials.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword)

	if err != nil {
		log.Println("Erreur QueryRow:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}


	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		log.Println("Erreur bcrypt:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Erreur génération token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
