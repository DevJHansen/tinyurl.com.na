package views

import (
	"net/http"

	"github.com/DevJHansen/tinyurl.com.na/pkg/models"
)

func (repo *Repository) Signup(w http.ResponseWriter, r *http.Request, user *models.User) {
	if user.ID != "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	err := repo.App.TemplateCache["signup.page.html"].Execute(w, "")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
