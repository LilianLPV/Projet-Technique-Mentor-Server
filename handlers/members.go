package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/config"
	"server/models"

	"github.com/go-chi/chi/v5"
)

func AddMemberHandler(w http.ResponseWriter, r *http.Request) {
	var member models.User
	err := json.NewDecoder(r.Body).Decode(&member)
	projectID := chi.URLParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	err = config.DB.QueryRow("SELECT id_user FROM users WHERE email = ?", member.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}
	_, err = config.DB.Exec("INSERT OR IGNORE INTO project_members (fk_project_id, fk_user_id, status) VALUES (?, ?, 'pending')",
		projectID,
		user.ID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Membre ajouté avec succès"))
}

func GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	rows, err := config.DB.Query(`
	SELECT u.id_user, u.username, u.email, 'accepted' as status
	FROM users u
	JOIN projects p ON p.fk_owner = u.id_user
	WHERE p.id_project = ?
	UNION
	SELECT u.id_user, u.username, u.email, pm.status 
	FROM users u 
	INNER JOIN project_members pm ON u.id_user = pm.fk_user_id 
	WHERE pm.fk_project_id = ?`, projectID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MemberWithStatus struct {
		ID       int    `json:"id_user"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Status   string `json:"status"`
	}

	var members []MemberWithStatus
	for rows.Next() {
		var member MemberWithStatus
		err := rows.Scan(&member.ID, &member.Username, &member.Email, &member.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		members = append(members, member)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func AssignTaskHandler(w http.ResponseWriter, r *http.Request) {
	var member models.User
	err := json.NewDecoder(r.Body).Decode(&member)
	taskID := chi.URLParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	err = config.DB.QueryRow("SELECT id_user FROM users WHERE email = ?", member.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}
	_, err = config.DB.Exec("INSERT OR IGNORE INTO task_members (fk_task, fk_user) VALUES (?, ?)", taskID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tâche assignée avec succès"))
}

func RemoveTaskMemberHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveTaskMemberHandler appelé")
	var member models.User
	err := json.NewDecoder(r.Body).Decode(&member)
	taskID := chi.URLParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	err = config.DB.QueryRow("SELECT id_user FROM users WHERE email = ?", member.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}
	result, err := config.DB.Exec("DELETE FROM task_members WHERE fk_task = ? AND fk_user = ?", taskID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	log.Println("Lignes supprimées:", rows)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Membre retiré avec succès"))
}

func LeaveProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	projectID := chi.URLParam(r, "id")
	_, err := config.DB.Exec("DELETE FROM project_members WHERE fk_project_id = ? AND fk_user_id = ?", projectID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Projet quitté avec succès"))
}

func RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	var member models.User
	err := json.NewDecoder(r.Body).Decode(&member)
	projectID := chi.URLParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	err = config.DB.QueryRow("SELECT id_user FROM users WHERE email = ?", member.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}
	_, err = config.DB.Exec("DELETE FROM project_members WHERE fk_project_id = ? AND fk_user_id = ?", projectID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Membre retiré avec succès"))
}

func GetInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)

	type Invitation struct {
		ProjectID    int    `json:"project_id"`
		ProjectTitle string `json:"project_title"`
	}

	rows, err := config.DB.Query(`
        SELECT p.id_project, p.title 
        FROM projects p
        JOIN project_members pm ON p.id_project = pm.fk_project_id
        WHERE pm.fk_user_id = ? AND pm.status = 'pending'`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var invitations []Invitation
	for rows.Next() {
		var inv Invitation
		rows.Scan(&inv.ProjectID, &inv.ProjectTitle)
		invitations = append(invitations, inv)
	}
	if invitations == nil {
		invitations = []Invitation{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitations)
}

func RespondInvitationHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	projectID := chi.URLParam(r, "id")

	var body struct {
		Accept bool `json:"accept"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if body.Accept {
		_, err := config.DB.Exec("UPDATE project_members SET status = 'accepted' WHERE fk_project_id = ? AND fk_user_id = ?", projectID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Invitation acceptée"))
	} else {
		_, err := config.DB.Exec("DELETE FROM project_members WHERE fk_project_id = ? AND fk_user_id = ?", projectID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Invitation refusée"))
	}
}
