package models

type Project struct {
    ID            int    `json:"id_project"`
    Title         string `json:"title"`
    Description   string `json:"description"`
    OwnerID       int    `json:"fk_owner"`
    OwnerUsername string `json:"owner_username"`
}