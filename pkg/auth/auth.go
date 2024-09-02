package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DevJHansen/tinyurl.com.na/pkg/config"
	"github.com/supabase-community/gotrue-go/types"
)

var Repo *Repository

type Repository struct {
	App *config.Config
}

func NewRepo(a *config.Config) {
	Repo = &Repository{
		App: a,
	}
}

func (repo *Repository) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body LoginBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Print("Error decoding body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := repo.App.Supabase.Auth.SignInWithEmailPassword(body.Email, body.Password)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    res.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (repo *Repository) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (repo *Repository) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var body SignupBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Print("Error decoding body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signupRequest := types.SignupRequest{
		Email:    body.Email,
		Password: body.Password,
		Phone:    "",
		Data:     map[string]interface{}{},
	}

	res, err := repo.App.Supabase.Auth.Signup(signupRequest)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, err = repo.App.Supabase.From("users").Insert(map[string]interface{}{
		"first_name": body.FirstName,
		"surname":    body.Surname,
		"uid":        res.User.ID,
	}, false, "", "representation", "").Execute()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    res.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
