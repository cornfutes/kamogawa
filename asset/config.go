package asset

import (
	_ "embed"
	"kamogawa/config"
	"kamogawa/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tdewolff/minify"
	minifyCss "github.com/tdewolff/minify/css"
)

var (
	//go:embed media/style.css
	styleCss    []byte
	styleCssMin []byte
	//go:embed media/cloud_logo_aws.png
	pngAWS []byte
	//go:embed media/cloud_logo_gcp.png
	pngGCP []byte
	//go:embed media/cloud_logo_azure.png
	pngAzure []byte
	//go:embed media/graphql.png
	pngGraphQl []byte
	//go:embed media/blog_traffic.png
	pngBlog1 []byte
	//go:embed media/blog_search.png
	pngBlog2 []byte
	//go:embed media/blog_splash.png
	pngBlog3 []byte
	//go:embed media/blog_login_error.png
	pngBlog4 []byte
	//go:embed media/blog_simple.png
	pngBlog5 []byte
	//go:embed media/blog_widget.png
	pngBlog6 []byte
	//go:embed media/blog_docker.png
	pngBlog7 []byte

	//go:embed media/splash.gif
	splash []byte
	//go:embed media/fuji.gif
	fuji []byte
	//go:embed media/ship.gif
	ship []byte
	//go:embed media/console.svg
	console []byte
	//go:embed media/phone.svg
	phone []byte
	//go:embed media/release.txt
	Release []byte
	//go:embed media/security.txt
	security []byte
	//go:embed media/legal.txt
	legal []byte
	//go:embed media/about.txt
	about []byte
	//go:embed media/api.txt
	api []byte

	//go:embed media/consent.png
	consent []byte

	staticHtml = map[string]string{
		"/":        "landing.html",
		"/docs":    "tbd.html",
		"/mission": "mission.html",
	}

	etag string
)

func init() {
	m := minify.New()
	preparecss(m)

	id := uuid.New()
	etag = id.String()
}

// Wire up HTML, css and media assets
func Config(r *gin.Engine) {
	r.HTMLRender = ConfigureHTMLRenderer()

	// Register static assets.
	if config.Env == config.Dev {
		r.StaticFile("/style.css", "asset/media/style.css")
	} else {
		r.GET("style.css", css(styleCss))
	}
	r.GET("cloud_logo_aws.png", png(pngAWS))
	r.GET("cloud_logo_gcp.png", png(pngGCP))
	r.GET("cloud_logo_azure.png", png(pngAzure))
	r.GET("graphql.png", png(pngGraphQl))
	r.GET("blog_traffic.png", png(pngBlog1))
	r.GET("blog_search.png", png(pngBlog2))
	r.GET("blog_splash.png", png(pngBlog3))
	r.GET("blog_login_error.png", png(pngBlog4))
	r.GET("blog_simple.png", png(pngBlog5))
	r.GET("blog_widget.png", png(pngBlog6))
	r.GET("blog_docker.png", png(pngBlog7))

	r.GET("splash.gif", gif(splash))
	r.GET("fuji.gif", gif(fuji))
	r.GET("ship.gif", gif(ship))
	r.GET("console.svg", svg(console))
	r.GET("phone.svg", svg(phone))
	r.GET("consent.png", png(consent))
	r.GET("legal.txt", TXT(legal))
	r.GET("security.txt", TXT(security))
	r.GET("about.txt", TXT(about))
	r.GET("api.txt", TXT(api))

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
	if err != nil {
		panic(err)
	}
}

func data(mime string, contents []byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		a := c.Request.Header["If-None-Match"]
		if len(a) > 0 && a[0] == etag {
			c.Data(http.StatusNotModified, mime, []byte{})
		} else {
			c.Header("ETag", etag)
			c.Header("ETag", etag)
			c.Data(http.StatusOK, mime, contents)
		}
	}
}

func css(contents []byte) func(c *gin.Context) {
	return data("text/css; charset=utf-8", contents)
}

func png(contents []byte) func(c *gin.Context) {
	return data("image/png", contents)
}

func gif(contents []byte) func(c *gin.Context) {
	return data("image/gif", contents)
}

func svg(contents []byte) func(c *gin.Context) {
	return data("image/svg+xml", contents)
}

func TXT(contents []byte) func(c *gin.Context) {
	return data("text/plain; charset=utf-8", contents)
}
