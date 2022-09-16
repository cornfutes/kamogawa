package asset

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

const baseLayoutPath = "view/layout/base.tmpl"
const unauthedLayoutsPath = "view/layout/unauthed/*.html"
const unauthedviewsGlob = "view/unauthed/*.html"
const authedLayoutsPath = "view/layout/authed/*.html"
const authedviewsGlob = "view/authed/*.html"

//go:embed view/*
var views embed.FS

// https://www.josephspurrier.com/how-to-embed-assets-in-go-1-16
/** Middleware to render HTML using templates. */
func ConfigureHTMLRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	unauthedviews, err := fs.Glob(views, unauthedviewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, unauthedView := range unauthedviews {
		t, err := template.ParseFS(views, baseLayoutPath, unauthedLayoutsPath, unauthedView)
		if err != nil {
			log.Fatal(err)
		}
		r.Add(filepath.Base(unauthedView), t)
	}

	authedviews, err := fs.Glob(views, authedviewsGlob)
	if err != nil {
		panic(err.Error())
	}
	for _, authedView := range authedviews {
		t, err := template.ParseFS(views, baseLayoutPath, authedLayoutsPath, authedView)
		if err != nil {
			log.Fatal(err)
		}
		r.Add(filepath.Base(authedView), t)
	}

	return r
}
