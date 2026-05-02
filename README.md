# Projet-Technique-Mentor-Server

API REST backend de l'application de gestion de projet, développée en Go dans le cadre d'un test technique pour le poste de mentor en filière Informatique & Cybersécurité.

## Technologies

- **Go 1.25**
- **Chi v5** — Router HTTP
- **JWT** — Authentification stateless
- **bcrypt** — Hashage des mots de passe
- **SQLite** — Base de données relationnelle locale

## Prérequis

- [Go 1.20+](https://go.dev/dl/)

## Installation

```bash
# Cloner le dépôt
git clone https://github.com/votre-compte/Projet-Technique-Mentor-Server.git
cd Projet-Technique-Mentor-Server

# Installer les dépendances
go mod tidy

# Lancer le serveur
go run main.go
```

Le serveur démarre sur **http://localhost:8080**

La base de données `projet.db` est créée automatiquement au premier démarrage.

## Structure

```
├── main.go          # Point d'entrée
├── config/
│   └── db.go        # Connexion SQLite + création des tables
├── models/          # Structures de données (User, Project, Task)
├── handlers/        # Logique métier (auth, projects, tasks, members)
├── middleware/      # Validation JWT
└── routes/          # Déclaration des routes + CORS
```

## Routes API

| Méthode | Route | Description |
|---------|-------|-------------|
| POST | /register | Créer un compte |
| POST | /login | Se connecter |
| GET | /projects | Liste des projets |
| POST | /projects | Créer un projet |
| GET | /projects/{id} | Détail d'un projet |
| PUT | /projects/{id} | Modifier un projet |
| DELETE | /projects/{id} | Supprimer un projet |
| GET | /projects/{id}/tasks | Liste des tâches |
| POST | /projects/{id}/tasks | Créer une tâche |
| PUT | /tasks/{id} | Modifier une tâche |
| DELETE | /tasks/{id} | Supprimer une tâche |
| POST | /projects/{id}/members | Inviter un membre |
| GET | /projects/{id}/members | Liste des membres |
| DELETE | /projects/{id}/members | Retirer un membre |
| DELETE | /projects/{id}/members/leave | Quitter un projet |
| PUT | /tasks/{id}/assign | Assigner un membre à une tâche |
| DELETE | /tasks/{id}/members | Retirer un membre d'une tâche |
| GET | /invitations | Voir ses invitations |
| PUT | /projects/{id}/invitations | Accepter/refuser une invitation |

Toutes les routes sauf `/register` et `/login` nécessitent un header :
```
Authorization: Bearer <token>
```

## Auteur

Lilian Le Piver