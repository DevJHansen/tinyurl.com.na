package render

import (
	"html/template"
	"path/filepath"
)

func CreateTemplateCache() (map[string]*template.Template, error) {
	var newCache = map[string]*template.Template{}

	pages, err := filepath.Glob("./static/*.page.html")

	if err != nil {
		return newCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)

		if err != nil {
			return newCache, err
		}

		newCache[name] = templateSet
	}

	return newCache, nil
}
