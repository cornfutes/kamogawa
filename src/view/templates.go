package view

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"

	"kamogawa/config"

	"github.com/gin-contrib/multitemplate"
)

//go:embed theme/*
var Views embed.FS

type Theme int8

const (
	Requiem Theme = iota
	Kubrick
)

func IsValidTheme(theme string) bool {
	if theme == Requiem.String() || theme == Kubrick.String() {
		return true
	}
	return false
}

func (s Theme) String() string {
	switch s {
	case Kubrick:
		return "kubrick"
	}
	return "requiem"
}

// https://www.josephspurrier.com/how-to-embed-assets-in-go-1-16
/** Middleware to render HTML using templates. */
func HTMLRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	registerTheme(r, Requiem)
	registerTheme(r, Kubrick)

	return r
}

func registerTheme(r multitemplate.Renderer, theme Theme) {
	baseLayoutPath := "theme/" + theme.String() + "/layout/*.tmpl"
	unauthedLayoutsPath := "theme/" + theme.String() + "/layout/unauthed/*.tmpl"
	authedLayoutsPath := "theme/" + theme.String() + "/layout/authed/*.tmpl"
	unauthedViewsGlob := "theme/" + theme.String() + "/unauthed/*.tmpl"
	authedViewsGlob := "theme/" + theme.String() + "/authed/*.tmpl"

	unauthedViews, err := fs.Glob(Views, unauthedViewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, unauthedView := range unauthedViews {
		t, err := template.ParseFS(Views, baseLayoutPath, unauthedLayoutsPath, unauthedView)
		if err != nil {
			log.Fatal(err)
		}
		if theme.String() == config.DefaultTheme {
			r.Add(filepath.Base(unauthedView), t)
		}
		r.Add(filepath.Base(unauthedView)+theme.String(), t)
	}

	authedViews, err := fs.Glob(Views, authedViewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, authedView := range authedViews {
		t, err := template.ParseFS(Views, baseLayoutPath, authedLayoutsPath, authedView)
		if err != nil {
			log.Fatal(err)
		}
		// Suffix theme name to end of template name
		if theme.String() == config.DefaultTheme {
			r.Add(filepath.Base(authedView), t)
		}
		r.Add(filepath.Base(authedView)+theme.String(), t)
	}
}
