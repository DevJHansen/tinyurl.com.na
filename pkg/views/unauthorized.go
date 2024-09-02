package views

import (
	"net/http"
)

func (repo *Repository) Unauthorized(w http.ResponseWriter, r *http.Request) {

	err := repo.App.TemplateCache["unauthorized.page.html"].Execute(w, "")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
