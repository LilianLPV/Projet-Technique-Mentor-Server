package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const JwtSecret = "ma_clé_secrète_pour_le_jwt"

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./projet.db")
	if err != nil {
		log.Fatal(err)
	}
	InitTables()
}

func InitTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id_user INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS projects (
			id_project INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,	
			description TEXT NOT NULL,
			fk_owner INTEGER NOT NULL,
			FOREIGN KEY (fk_owner) REFERENCES users(id_user)
		);
		CREATE TABLE IF NOT EXISTS tasks (
			id_task INTEGER PRIMARY KEY AUTOINCREMENT,
			fk_project INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			status TEXT NOT NULL,
			FOREIGN KEY (fk_project) REFERENCES projects(id_project)
		);
		CREATE TABLE IF NOT EXISTS task_members (
			fk_task INTEGER NOT NULL,
			fk_user INTEGER NOT NULL,
			PRIMARY KEY (fk_task, fk_user),
			FOREIGN KEY (fk_task) REFERENCES tasks(id_task),
			FOREIGN KEY (fk_user) REFERENCES users(id_user)
		);
		CREATE TABLE IF NOT EXISTS project_members (
			fk_project_id INTEGER,
			fk_user_id INTEGER,
			status TEXT DEFAULT 'pending',
			PRIMARY KEY (fk_project_id, fk_user_id),
			FOREIGN KEY (fk_project_id) REFERENCES projects(id_project),
			FOREIGN KEY (fk_user_id) REFERENCES users(id_user)
		);
		`)
	if err != nil {
		log.Fatal(err)
	}
}
