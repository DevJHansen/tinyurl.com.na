package views

import (
	"net/http"

	"github.com/DevJHansen/tinyurl.com.na/pkg/models"
)

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request, user *models.User) {

	data := map[string]bool{
		"LoggedIn": false,
	}

	if user.ID != "" {
		data["LoggedIn"] = true
	}

	err := repo.App.TemplateCache["home.page.html"].Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
