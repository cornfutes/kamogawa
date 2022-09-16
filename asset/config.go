package asset

import (
	_ "embed"
	"kamogawa/config"
	"kamogawa/core"
	"kamogawa/view"

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
	//go:embed media/screensaver.mp4
	mp4Screensaver []byte
	//go:embed media/cloud_logo_aws.png
	pngAWS []byte
	//go:embed media/cloud_logo_gcp.png
	pngGCP []byte
	//go:embed media/cloud_logo_azure.png
	pngAzure []byte
	//go:embed media/graphql.png
	pngGraphQl []byte
	//go:embed media/blog_traffic.gif
	pngBlog1 []byte
	//go:embed media/blog_search.gif
	pngBlog2 []byte
	//go:embed media/blog_splash.jpg
	pngBlog3 []byte
	//go:embed media/blog_login_error.gif
	pngBlog4 []byte
	//go:embed media/blog_simple.gif
	pngBlog5 []byte
	//go:embed media/blog_widget.gif
	pngBlog6 []byte
	//go:embed media/blog_docker.gif
	pngBlog7 []byte

	//go:embed media/splash_landing.gif
	splashLanding []byte
	//go:embed media/splash_fuji.gif
	splashFuji []byte
	//go:embed media/splash_ship.gif
	splashShip []byte
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
	//go:embed media/nft.gif
	gifProfile []byte
	//go:embed media/big_sur.jpg
	jpgBigSur []byte

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
	// TODO: if people hit different web servers, this would be different tags
	etag = id.String()
}

// Wire up HTML, css and media assets
func Config(r *gin.Engine) {
	r.HTMLRender = view.HTMLRenderer()

	// Register static assets.
	if config.Env == config.Dev {
		r.StaticFile("/style.css", "asset/media/style.css")
		r.StaticFile("/style_glass.css", "asset/media/style_glass.css")
	} else {
		r.GET("style.css", css(styleCssMin))
	}
	r.GET("screensaver.mp4", mp4(mp4Screensaver))
	r.GET("cloud_logo_aws.png", png(pngAWS))
	r.GET("cloud_logo_gcp.png", png(pngGCP))
	r.GET("cloud_logo_azure.png", png(pngAzure))
	r.GET("graphql.png", png(pngGraphQl))
	r.GET("blog_traffic.gif", gif(pngBlog1))
	r.GET("blog_search.gif", gif(pngBlog2))
	r.GET("blog_splash.jpg", jpg(pngBlog3))
	r.GET("blog_login_error.gif", gif(pngBlog4))
	r.GET("blog_simple.gif", png(pngBlog5))
	r.GET("blog_widget.gif", png(pngBlog6))
	r.GET("blog_docker.gif", gif(pngBlog7))

	r.GET("splash_landing.gif", gif(splashLanding))
	r.GET("splash_fuji.gif", gif(splashFuji))
	r.GET("splash_ship.gif", gif(splashShip))
	r.GET("console.svg", svg(console))
	r.GET("phone.svg", svg(phone))
	r.GET("consent.png", png(consent))
	r.GET("legal.txt", TXT(legal))
	r.GET("security.txt", TXT(security))
	r.GET("about.txt", TXT(about))
	r.GET("api.txt", TXT(api))
	r.GET("nft.gif", gif(gifProfile))
	r.GET("big_sur.jpg", jpg(jpgBigSur))

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
			c.Data(http.StatusOK, mime, contents)
		}
	}
}

func css(contents []byte) func(c *gin.Context) {
	// TODO: by default browsers utf-8
	return data("text/css; charset=utf-8", contents)
}

func png(contents []byte) func(c *gin.Context) {
	return data("image/png", contents)
}

func jpg(contents []byte) func(c *gin.Context) {
	return data("image/jpg", contents)
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

func mp4(contents []byte) func(c *gin.Context) {
	return data("video/mp4", contents)
}
