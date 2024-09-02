package routes

import (
	"net/http"

	"github.com/DevJHansen/tinyurl.com.na/pkg/auth"
	"github.com/DevJHansen/tinyurl.com.na/pkg/redirect"
	"github.com/DevJHansen/tinyurl.com.na/pkg/views"
)

func Routes(router *http.ServeMux) {

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, views.Repo.Home)
	})

	router.HandleFunc("GET /dashboard", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, views.Repo.Dashboard)
	})

	router.HandleFunc("GET /signup", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, views.Repo.Signup)
	})

	router.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.HandleSignup(w, r)
	})

	router.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, views.Repo.Login)
	})

	router.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.HandleLogin(w, r)
	})

	router.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.HandleLogout(w, r)
	})

	router.HandleFunc("GET /{key}", func(w http.ResponseWriter, r *http.Request) {
		redirect.Repo.HandleRedirect(w, r)
	})

	router.HandleFunc("POST /redirects", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, redirect.Repo.HandleCreateRedirect)
	})

	router.HandleFunc("GET /redirects", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.CheckAuth(w, r, redirect.Repo.HandleGetRedirects)
	})

	router.HandleFunc("DELETE /{key}", func(w http.ResponseWriter, r *http.Request) {
		auth.Repo.RequireAuth(w, r, redirect.Repo.HandleDeleteRedirect)
	})
}
