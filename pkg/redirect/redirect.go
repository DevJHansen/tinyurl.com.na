package redirect

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DevJHansen/tinyurl.com.na/pkg/config"
	"github.com/DevJHansen/tinyurl.com.na/pkg/models"
	"github.com/supabase-community/postgrest-go"
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

func (repo *Repository) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	target, ok := (*repo.App.RedirectCache)[key]

	if !ok {
		keyLookUpRes, _, _ := repo.App.Supabase.From("redirects").Select("*", "exact", false).Eq("key", key).Neq("deleted", "true").Execute()

		var redirectsWithKey []models.Redirect
		err := json.Unmarshal(keyLookUpRes, &redirectsWithKey)

		if err != nil {
			log.Print("Error during redirects unmarshal")
			http.Error(w, "Error redirecting to target", http.StatusInternalServerError)
			return
		}

		if len(redirectsWithKey) == 0 {
			http.Error(w, "Target not found", http.StatusNotFound)
			return
		}

		go LogRedirectAnalytics(r, key, repo.App)

		target = redirectsWithKey[0].Target
		(*repo.App.RedirectCache)[key] = target
		http.Redirect(w, r, target, http.StatusFound)
		return
	}

	go LogRedirectAnalytics(r, key, repo.App)

	http.Redirect(w, r, target, http.StatusFound)
}

func (repo *Repository) HandleCreateRedirect(w http.ResponseWriter, r *http.Request, user *models.User) {
	var body models.NewRedirectReqBody
	owner := "system"

	if user.ID != "" {
		owner = user.ID
	}

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Print("Error decoding body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	target, err := ProcessUrl(body.Target)

	if err != nil {
		log.Print("Error processing target")
		http.Error(w, "Please provide a target url to shorten.", http.StatusBadRequest)
		return
	}

	key, err := GenerateUID()

	if err != nil {
		log.Print("Error generating uid")
		http.Error(w, fmt.Sprintf("Error generating mini link for: %q", target), http.StatusInternalServerError)
		return
	}

	keyLookUpRes, _, _ := repo.App.Supabase.From("redirects").Select("*", "exact", false).Eq("key", key).Execute()

	var redirectsWithKey []models.Redirect
	err = json.Unmarshal(keyLookUpRes, &redirectsWithKey)

	if err != nil {
		log.Print("Error during redirects unmarshal")
		http.Error(w, fmt.Sprintf("Error generating mini link for: %q. Please try again.", target), http.StatusInternalServerError)
		return
	}

	if len(redirectsWithKey) > 0 {
		http.Error(w, fmt.Sprintf("Error generating mini link for: %q. Please try again.", target), http.StatusInternalServerError)
		return
	}

	res, _, err := repo.App.Supabase.From("redirects").Insert(map[string]interface{}{
		"key":    key,
		"target": target,
		"owner":  owner,
	}, false, "", "representation", "").Execute()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating mini link for: %q. Please try again.", target), http.StatusInternalServerError)
		return
	}

	var response []models.Redirect
	err = json.Unmarshal(res, &response)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating mini link for: %q. Please try again.", target), http.StatusInternalServerError)
		log.Print("Error during new redirect unmarshal")
		return
	}

	(*repo.App.RedirectCache)[key] = target

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response[0]); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (repo *Repository) HandleDeleteRedirect(w http.ResponseWriter, r *http.Request, user *models.User) {
	key := r.PathValue("key")

	updateData := map[string]interface{}{
		"deleted": true,
	}

	_, _, err := repo.App.Supabase.From("redirects").Update(updateData, "", "").Eq("key", key).Eq("owner", user.ID).Execute()

	if err != nil {
		http.Error(w, "Error deleting target", http.StatusInternalServerError)
		return
	}

	_, ok := (*repo.App.RedirectCache)[key]

	if ok {
		delete((*repo.App.RedirectCache), key)
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{})
}

func (repo *Repository) HandleGetRedirects(w http.ResponseWriter, r *http.Request, user *models.User) {
	page := r.URL.Query().Get("page")
	var redirects []models.Redirect

	if page == "" {
		page = "1"
	}

	num, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Error parsing page number", http.StatusBadRequest)
		return
	}

	startIndex := ((num - 1) * 20)
	endIndex := startIndex + 19

	opts := postgrest.OrderOpts{Ascending: true, NullsFirst: false, ForeignTable: ""}
	res, count, err := repo.App.Supabase.From("redirects").Select("*", "exact", false).Eq("owner", user.ID).Neq("deleted", "true").Range(startIndex, endIndex, "").Order("created_at", &opts).Execute()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(res, &redirects)

	if err != nil {
		log.Print("Error during redirects unmarshal")
		http.Error(w, "Error fetching mini links", http.StatusInternalServerError)
		return
	}

	var resBody = models.GetRedirectsRes{
		Page:  num,
		Data:  redirects,
		Count: count,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}
