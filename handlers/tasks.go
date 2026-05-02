package handlers

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/models"

	"github.com/go-chi/chi/v5"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	projectID := chi.URLParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = config.DB.Exec("INSERT INTO tasks (title, description, fk_project, status) VALUES (?, ?, ?, ?)",
		task.Title,
		task.Description,
		projectID,
		task.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Tâche créée avec succès"))
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	rows, err := config.DB.Query(`
    SELECT t.id_task, t.title, t.description, t.status, t.fk_project
    FROM tasks t
    WHERE t.fk_project = ?`, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.ProjectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	for i, task := range tasks {
		rows2, err := config.DB.Query(`
    SELECT u.email FROM users u
    JOIN task_members tm ON u.id_user = tm.fk_user
    WHERE tm.fk_task = ?`, task.ID)
		if err != nil {
			continue
		}
		defer rows2.Close()
		var members []string
		for rows2.Next() {
			var username string
			rows2.Scan(&username)
			members = append(members, username)
		}
		tasks[i].AssignedTo = members
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskID := chi.URLParam(r, "id")
result, err := config.DB.Exec(`
    UPDATE tasks SET title = ?, description = ?, status = ? 
    WHERE id_task = ? 
    AND fk_project IN (
        SELECT id_project FROM projects WHERE fk_owner = ?
        UNION
        SELECT fk_project_id FROM project_members WHERE fk_user_id = ?
    )`,
		task.Title,
		task.Description,
		task.Status,
		taskID,
		userID,
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
		http.Error(w, "Tâche non trouvée ou pas autorisée", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tâche modifiée avec succès"))
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(float64)
	taskID := chi.URLParam(r, "id")
	result, err := config.DB.Exec("DELETE FROM tasks WHERE id_task = ? AND fk_project IN (SELECT id_project FROM projects WHERE fk_owner = ?)", taskID, userID)
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
		http.Error(w, "Tâche non trouvée ou pas autorisée", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tâche supprimée avec succès"))
}
