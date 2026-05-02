package models

type Task struct {
    ID          int     `json:"id_task"`
    UserID      *int    `json:"user_id"`
    ProjectID   int     `json:"project_id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Status      string  `json:"status"`
    AssignedTo []string `json:"assigned_to"`

}