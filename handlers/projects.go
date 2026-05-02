package handlers

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/models"

	"github.com/go-chi/chi/v5"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var projects models.Project
	err := json.NewDecoder(r.Body).Decode(&projects)
	userID := r.Context().Value("user_id").(float64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = config.DB.Exec("INSERT INTO projects (title, description, fk_owner) VALUES (?, ?, ?)",
		projects.Title,
		projects.Description,
		userID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Projet créé avec succès"))
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	rows, err := config.DB.Query(`
    SELECT id_project, title, description, fk_owner FROM projects WHERE fk_owner = ?
    UNION
    SELECT p.id_project, p.title, p.description, p.fk_owner 
    FROM projects p 
    JOIN project_members pm ON p.id_project = pm.fk_project_id 
    WHERE pm.fk_user_id = ? AND pm.status = 'accepted'`, userID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.OwnerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}
	if projects == nil {
		projects = []models.Project{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	projectID := chi.URLParam(r, "id")
	result, err := config.DB.Exec("UPDATE projects SET title = ?, description = ? WHERE id_project = ? AND fk_owner = ?",
		project.Title,
		project.Description,
		projectID,
		userID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Projet non trouvé ou pas autorisé", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Projet modifié avec succès"))
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	projectID := chi.URLParam(r, "id")
	result, err := config.DB.Exec("DELETE FROM projects WHERE id_project = ? AND fk_owner = ?", projectID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Projet non trouvé ou pas autorisé", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Projet supprimé avec succès"))
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	projectID := chi.URLParam(r, "id")
	var project models.Project
	err := config.DB.QueryRow(`
    SELECT p.id_project, p.title, p.description, p.fk_owner, u.username 
    FROM projects p
    JOIN users u ON p.fk_owner = u.id_user
    WHERE p.id_project = ? AND (p.fk_owner = ? OR p.id_project IN (SELECT fk_project_id FROM project_members WHERE fk_user_id = ?))`,
		projectID, userID, userID).Scan(&project.ID, &project.Title, &project.Description, &project.OwnerID, &project.OwnerUsername)
	if err != nil {
		http.Error(w, "Projet non trouvé", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
