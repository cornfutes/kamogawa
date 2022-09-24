package media

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"kamogawa/config"
	"kamogawa/core"
	"kamogawa/view"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tdewolff/minify"
	minifyCss "github.com/tdewolff/minify/css"
)

type MimeType string

const (
	MimeTypeMP4 MimeType = "video/mp4"
	MimeTypeTXT MimeType = "text/plain; charset=utf-8"
	MimeTypePNG MimeType = "image/png"
	MimeTypeCSS MimeType = "text/css; charset=utf-8"
	MimeTypeJPG MimeType = "image/jpg"
	MimeTypeGIF MimeType = "image/gif"
	MimeTypeSVG MimeType = "image/svg+xml"
)

var (
	//go:embed asset/style.css
	styleCss    []byte
	styleCssMin []byte
	//go:embed asset/style_glass.css
	styleCssKubrick            []byte
	styleCssStyleCssKubrickMin []byte

	//go:embed asset/screensaver.mp4
	mp4Screensaver []byte
	//go:embed asset/cloud_logo_aws.png
	pngAWS []byte
	//go:embed asset/cloud_logo_gcp.png
	pngGCP []byte
	//go:embed asset/cloud_logo_azure.png
	pngAzure []byte

	//go:embed asset/blog_traffic.gif
	pngBlog1 []byte
	//go:embed asset/blog_search.gif
	pngBlog2 []byte
	//go:embed asset/blog_splash.jpg
	pngBlog3 []byte
	//go:embed asset/blog_login_error.gif
	pngBlog4 []byte
	//go:embed asset/blog_simple.gif
	pngBlog5 []byte
	//go:embed asset/blog_widget.gif
	pngBlog6 []byte
	//go:embed asset/blog_docker.gif
	pngBlog7 []byte

	//go:embed asset/splash_landing.gif
	splashLanding []byte
	//go:embed asset/splash_landing_hd.gif
	splashLandingHD []byte
	//go:embed asset/landing_clock.png
	landingClockPng []byte
	//go:embed asset/landing_screenshot.png
	landingScreenshot []byte

	//go:embed asset/splash_fuji.gif
	splashFuji []byte
	//go:embed asset/splash_ship.gif
	splashShip []byte
	//go:embed asset/console.svg
	console []byte
	//go:embed asset/phone.svg
	phone []byte
	//go:embed asset/release.txt
	Release []byte
	//go:embed asset/security.txt
	security []byte
	//go:embed asset/legal.txt
	legal []byte
	//go:embed asset/about.txt
	about []byte
	//go:embed asset/api.txt
	api []byte
	//go:embed asset/nft.gif
	gifProfile []byte
	//go:embed asset/big_sur.jpg
	jpgBigSur []byte

	//go:embed asset/consent.png
	consent []byte

	staticHtml = map[string]string{
		// "/":        "landing.tmpl",
		"/docs":    "tbd.tmpl",
		"/mission": "mission.tmpl",
	}

	etag string
)

func init() {
	m := minify.New()
	preparecss(m)

	id := uuid.New()
	// TODO: if people hit different web servers, this would be different tags
	etag = id.String()
}

const MediaBaseUrl = "static/"

// Wire up HTML, css and media assets
func Register(r *gin.Engine) {
	x := view.HTMLRenderer()
	t, err := template.ParseFS(view.Views, "theme/requiem/unauthed/landing.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	x.Add("xanax", t)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "xanax", gin.H{})
	})
	r.HTMLRender = x

	// Register static assets.
	if config.Env == config.Dev {
		r.StaticFile("/style.css", "media/asset/style.css")
		r.StaticFile("/style_glass.css", "media/asset/style_glass.css")
	} else {
		r.GET("style.css", Data(MimeTypeCSS, styleCssMin))
		r.GET("style_glass.css", Data(MimeTypeCSS, styleCssStyleCssKubrickMin))
	}
	r.GET(MediaBaseUrl+"screensaver.mp4", Data(MimeTypeMP4, mp4Screensaver))
	r.GET(MediaBaseUrl+"cloud_logo_aws.png", Data(MimeTypePNG, pngAWS))
	r.GET(MediaBaseUrl+"cloud_logo_gcp.png", Data(MimeTypePNG, pngGCP))
	r.GET(MediaBaseUrl+"cloud_logo_azure.png", Data(MimeTypePNG, pngAzure))
	r.GET(MediaBaseUrl+"blog_traffic.gif", Data(MimeTypeGIF, pngBlog1))
	r.GET(MediaBaseUrl+"blog_search.gif", Data(MimeTypeGIF, pngBlog2))
	r.GET(MediaBaseUrl+"blog_splash.jpg", Data(MimeTypeJPG, pngBlog3))
	r.GET(MediaBaseUrl+"blog_login_error.gif", Data(MimeTypeGIF, pngBlog4))
	r.GET(MediaBaseUrl+"blog_simple.gif", Data(MimeTypePNG, pngBlog5))
	r.GET(MediaBaseUrl+"blog_widget.gif", Data(MimeTypePNG, pngBlog6))
	r.GET(MediaBaseUrl+"blog_docker.gif", Data(MimeTypeGIF, pngBlog7))

	r.GET(MediaBaseUrl+"splash_landing.gif", Data(MimeTypeGIF, splashLanding))
	r.GET(MediaBaseUrl+"splash_landing_hd.gif", Data(MimeTypeGIF, splashLandingHD))
	r.GET(MediaBaseUrl+"splash_fuji.gif", Data(MimeTypeGIF, splashFuji))
	r.GET(MediaBaseUrl+"splash_ship.gif", Data(MimeTypeGIF, splashShip))
	r.GET(MediaBaseUrl+"console.svg", Data(MimeTypeSVG, console))
	r.GET(MediaBaseUrl+"phone.svg", Data(MimeTypeSVG, phone))
	r.GET(MediaBaseUrl+"consent.png", Data(MimeTypePNG, consent))
	r.GET(MediaBaseUrl+"nft.gif", Data(MimeTypeGIF, gifProfile))
	r.GET(MediaBaseUrl+"big_sur.jpg", Data(MimeTypeJPG, jpgBigSur))
	r.GET(MediaBaseUrl+"landing_clock.png", Data(MimeTypePNG, landingClockPng))
	r.GET(MediaBaseUrl+"landing_screenshot.png", Data(MimeTypePNG, landingScreenshot))

	r.GET("legal.txt", Data(MimeTypeTXT, legal))
	r.GET("security.txt", Data(MimeTypeTXT, security))
	r.GET("about.txt", Data(MimeTypeTXT, about))
	r.GET("api.txt", Data(MimeTypeTXT, api))

	// Register static views.
	for route, file := range staticHtml {
		func(f string) {
			r.GET(route, func(c *gin.Context) {
				core.HTMLWithGlobalState(c, f, gin.H{})
			})
		}(file)
	}
}

func preparecss(m *minify.M) {
	m.AddFunc("text/css", minifyCss.Minify)
	var err error
	styleCssMin, err = m.Bytes("text/css", styleCss)
	styleCssStyleCssKubrickMin, err = m.Bytes("text/css", styleCssKubrick)
	if err != nil {
		panic(err)
	}
}

func Data(mime MimeType, contents []byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		a := c.Request.Header["If-None-Match"]
		if len(a) > 0 && a[0] == etag {
			c.Data(http.StatusNotModified, string(mime), []byte{})
		} else {
			c.Header("ETag", etag)
			c.Writer.Header().Set("Cache-Control", "public, max-age=1800, immutable")
			c.Data(http.StatusOK, string(mime), contents)
		}
	}
}
