package asset

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

const baseLayoutPath = "view/layout/matrix/*.tmpl"
const unauthedLayoutsPath = "view/layout/matrix/unauthed/*.html"
const authedLayoutsPath = "view/layout/matrix/authed/*.html"
const unauthedViewsGlob = "view/unauthed/*.html"
const authedViewsGlob = "view/authed/*.html"

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
