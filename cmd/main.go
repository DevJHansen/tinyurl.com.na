package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DevJHansen/tinyurl.com.na/pkg/auth"
	"github.com/DevJHansen/tinyurl.com.na/pkg/config"
	"github.com/DevJHansen/tinyurl.com.na/pkg/redirect"
	"github.com/DevJHansen/tinyurl.com.na/pkg/render"
	"github.com/DevJHansen/tinyurl.com.na/pkg/routes"
	"github.com/DevJHansen/tinyurl.com.na/pkg/views"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var app config.Config
var redirectCache = make(map[string]string)
var envs = make(map[string]string)

func main() {
	router := http.NewServeMux()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
		return
	}

	API_URL := os.Getenv("DB_API_URL")
	API_KEY := os.Getenv("DB_API_KEY")
	SUPABASE_JWT_SECRET := os.Getenv("SUPABASE_JWT_SECRET")
	IP_INFO_API_KEY := os.Getenv("IP_INFO_API_KEY")

	envs["API_URL"] = API_URL
	envs["API_KEY"] = API_KEY
	envs["SUPABASE_JWT_SECRET"] = SUPABASE_JWT_SECRET
	envs["IP_INFO_API_KEY"] = IP_INFO_API_KEY

	client, err := supabase.NewClient(API_URL, API_KEY, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	templateCache, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return
	}

	app.Supabase = client
	app.RedirectCache = &redirectCache
	app.Envs = &envs
	app.TemplateCache = templateCache
	auth.NewRepo(&app)
	redirect.NewRepo(&app)
	views.NewRepo(&app)

	routes.Routes(router)

	http.ListenAndServe(":8080", router)
}
