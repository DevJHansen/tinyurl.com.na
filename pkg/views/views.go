package views

import "github.com/DevJHansen/tinyurl.com.na/pkg/config"

var Repo *Repository

type Repository struct {
	App *config.Config
}

func NewRepo(a *config.Config) {
	Repo = &Repository{
		App: a,
	}
}
