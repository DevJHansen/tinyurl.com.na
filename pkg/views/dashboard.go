package views

import (
	"net/http"

	"github.com/DevJHansen/tinyurl.com.na/pkg/models"
)

func (repo *Repository) Dashboard(w http.ResponseWriter, r *http.Request, user *models.User) {
	if user.ID == "" {
		repo.Unauthorized(w, r)
		return
	}

	err := repo.App.TemplateCache["dashboard.page.html"].Execute(w, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
