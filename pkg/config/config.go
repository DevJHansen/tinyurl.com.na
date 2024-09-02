package config

import (
	"html/template"

	"github.com/supabase-community/supabase-go"
)

type Config struct {
	Supabase      *supabase.Client
	RedirectCache *map[string]string
	Envs          *map[string]string
	TemplateCache map[string]*template.Template
}
