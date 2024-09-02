package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DevJHansen/tinyurl.com.na/pkg/models"
	"github.com/go-resty/resty/v2"
)

func (repo *Repository) RequireAuth(w http.ResponseWriter, r *http.Request, next models.AuthHandler) {
	cookie, err := r.Cookie("auth-token")

	if err == nil {
		user, err := repo.decodeJWT(cookie.Value)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if user.ID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r, &user)
		return
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := repo.decodeJWT(token)

	if user.ID != "" && err == nil {
		next(w, r, &user)
		return
	}

	user, err = repo.validateApiKey(token)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	if user.ID != "" {
		next(w, r, &user)
		return
	}

	http.Error(w, "Unauthorized", http.StatusInternalServerError)
}

func (repo *Repository) CheckAuth(w http.ResponseWriter, r *http.Request, next models.AuthHandler) {
	cookie, err := r.Cookie("auth-token")
	user := models.User{}

	if err == nil {
		user, err := repo.decodeJWT(cookie.Value)

		if err != nil {
			next(w, r, &user)
			return
		}

		next(w, r, &user)
		return
	}

	next(w, r, &user)
}

func (repo *Repository) validateApiKey(apiKey string) (models.User, error) {
	keyLookUpRes, _, err := repo.App.Supabase.From("api_keys").Select("*", "exact", false).Eq("api_key", apiKey).Execute()

	if err != nil {
		return models.User{}, err
	}

	var keys []models.ApiKey
	err = json.Unmarshal(keyLookUpRes, &keys)

	if err != nil {
		return models.User{}, err
	}

	if len(keys) == 0 {
		return models.User{}, fmt.Errorf("invalid api key")
	}

	user, err := repo.getUserByID(keys[0].Owner)

	if err != nil {
		return models.User{}, fmt.Errorf("error fetching user")
	}

	return user, nil
}

func (repo *Repository) decodeJWT(t string) (models.User, error) {
	url := "https://thodvwdwfsuwlfntzdoi.supabase.co/auth/v1/user"
	client := resty.New()

	var result map[string]interface{}
	var user models.User
	resp, err := client.R().
		SetHeader("apikey", (*repo.App.Envs)["API_KEY"]).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", t)).
		SetResult(&result).
		Get(url)

	if err != nil {
		return user, err
	}

	if resp.IsError() {
		return user, fmt.Errorf("error response code: %d", resp.StatusCode())
	}

	if msg, ok := result["message"].(string); ok {
		return user, fmt.Errorf("auth Error: %s", msg)
	}

	if id, ok := result["id"].(string); ok {
		user = models.User{
			ID:        id,
			Aud:       result["aud"].(string),
			Role:      result["role"].(string),
			Email:     result["email"].(string),
			CreatedAt: time.Time{},
		}

		if createdAtStr, ok := result["created_at"].(string); ok {
			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
			if err == nil {
				user.CreatedAt = createdAt
			}
		}

		return user, nil
	}

	return user, fmt.Errorf("unexpected response structure")
}

func (repo *Repository) getUserByID(id string) (models.User, error) {
	url := fmt.Sprintf("https://thodvwdwfsuwlfntzdoi.supabase.co/auth/v1/user/%s", id)
	apiKey := (*repo.App.Envs)["API_KEY"]
	client := resty.New()

	var result map[string]interface{}
	var user models.User
	resp, err := client.R().
		SetHeader("apikey", apiKey).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
		SetResult(&result).
		Get(url)

	if err != nil {
		return user, err
	}

	if resp.IsError() {
		return user, fmt.Errorf("error response code: %d", resp.StatusCode())
	}

	if msg, ok := result["message"].(string); ok {
		return user, fmt.Errorf("auth Error: %s", msg)
	}

	if id, ok := result["id"].(string); ok {
		user = models.User{
			ID:        id,
			Aud:       result["aud"].(string),
			Role:      result["role"].(string),
			Email:     result["email"].(string),
			CreatedAt: time.Time{},
		}

		if createdAtStr, ok := result["created_at"].(string); ok {
			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
			if err == nil {
				user.CreatedAt = createdAt
			}
		}

		return user, nil
	}

	return user, fmt.Errorf("unexpected response structure")
}
