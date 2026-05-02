package routes

import (
	"fmt"
	"net/http"
	"server/handlers"
	"server/middleware"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetupRoutes(router *chi.Mux) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/projects", handlers.CreateProjectHandler)
		r.Get("/projects", handlers.GetProjectsHandler)
		r.Put("/projects/{id}", handlers.UpdateProjectHandler)
		r.Delete("/projects/{id}", handlers.DeleteProjectHandler)
		r.Post("/projects/{id}/tasks", handlers.CreateTaskHandler)
		r.Get("/projects/{id}/tasks", handlers.GetTasksHandler)
		r.Put("/tasks/{id}", handlers.UpdateTaskHandler)
		r.Delete("/tasks/{id}", handlers.DeleteTaskHandler)
		r.Post("/projects/{id}/members", handlers.AddMemberHandler)
		r.Get("/projects/{id}/members", handlers.GetMembersHandler)
		r.Put("/tasks/{id}/assign", handlers.AssignTaskHandler)
		r.Delete("/tasks/{id}/members", handlers.RemoveTaskMemberHandler)
		r.Get("/projects/{id}", handlers.GetProjectHandler)
		r.Delete("/projects/{id}/members/leave", handlers.LeaveProjectHandler)
		r.Delete("/projects/{id}/members", handlers.RemoveProjectMemberHandler)
		r.Get("/invitations", handlers.GetInvitationsHandler)
		r.Put("/projects/{id}/invitations", handlers.RespondInvitationHandler)

	})
	router.Get("/test-bcrypt", func(w http.ResponseWriter, r *http.Request) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
		err := bcrypt.CompareHashAndPassword(hash, []byte("test"))
		if err != nil {
			w.Write([]byte("ERREUR: " + err.Error()))
		} else {
			w.Write([]byte("OK"))
		}
	})
	router.Post("/register", handlers.RegisterHandler)
	router.Post("/login", handlers.LoginHandler)

}
