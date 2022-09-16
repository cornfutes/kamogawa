package asset

import (
	"embed"
	"html/template"
	"io/fs"
	"kamogawa/config"
	"log"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

var baseLayoutPath = "view/layout/" + config.Theme + "/*.tmpl"
var unauthedLayoutsPath = "view/layout/" + config.Theme + "/unauthed/*.html"
var authedLayoutsPath = "view/layout/" + config.Theme + "/authed/*.html"
var unauthedViewsGlob = "view/unauthed/*.html"
var authedViewsGlob = "view/authed/*.html"

//go:embed view/*
var views embed.FS

// https://www.josephspurrier.com/how-to-embed-assets-in-go-1-16
/** Middleware to render HTML using templates. */
func ConfigureHTMLRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	unauthedViews, err := fs.Glob(views, unauthedViewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, unauthedView := range unauthedViews {
		t, err := template.ParseFS(views, baseLayoutPath, unauthedLayoutsPath, unauthedView)
		if err != nil {
			log.Fatal(err)
		}
		r.Add(filepath.Base(unauthedView), t)
	}

	authedViews, err := fs.Glob(views, authedViewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, authedView := range authedViews {
		t, err := template.ParseFS(views, baseLayoutPath, authedLayoutsPath, authedView)
		if err != nil {
			log.Fatal(err)
		}
		r.Add(filepath.Base(authedView), t)
	}

	return r
}
