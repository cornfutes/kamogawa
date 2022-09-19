package main

import (
	"kamogawa/cache/gcpcache/gcecache"
	"log"
	"net/http"
	"os"
	"strings"

	"kamogawa/config"
	"kamogawa/core"
	"kamogawa/handler"
	"kamogawa/identity"
	"kamogawa/media"
	"kamogawa/types"
	"kamogawa/view"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	shimogawaUrl := os.Getenv("SHIMOGAWA_URL")
	if len(shimogawaUrl) == 0 {
		log.Panic("$SHIMOGAWA_URL not set")
	}
	db = core.InitDB(shimogawaUrl)

	// TODO: remove. for prototyping purposes.
	// note: we currently ignore the stored DB password.
	for email, password := range identity.UsersInMemory {
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			panic("Error setting up initializing user credentials.")
		}
		db.FirstOrCreate(&types.User{
			Email:    email,
			Password: string(encryptedPassword),
		})
	}
}

func main() {
	if config.Env != config.Dev {
		gin.SetMode("release")
	}
	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(identity.SetAuthContext())
	media.Register(r)

	r.GET("/demo", handler.Demo)
	r.GET("/blog", handler.Blog)
	r.GET("/blog_showhn", handler.Blog1)

	r.GET("/login", handler.Login)
	r.GET("/reset", handler.Reset)
	r.POST("/login", identity.HandleLogin)
	r.POST("/reset", identity.HandleReset)
	r.GET("/loggedout", func(c *gin.Context) {
		core.HTMLWithGlobalState(c, "loggedout.tmpl", gin.H{})
	})

	r.GET("/set_theme", func(c *gin.Context) {
		theme := c.Query("theme")
		if !view.IsValidTheme(theme) {
			theme = config.DefaultTheme
		}
		c.SetCookie(identity.CookieKeyTheme, theme, 2147483647, "/", config.Host, config.CookieHttpsOnly, true)
		c.Redirect(http.StatusFound, c.Request.Referer())
	})

	authed := r.Group("/", identity.GateAuth())
	{
		authed.GET("/static/release.txt", media.Data(media.MimeTypeTXT, media.Release, false))

		authed.GET("/glass_sample", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "glass_sample.tmpl", gin.H{})
			c.Abort()
		})

		// ******** Begin left nav
		authed.GET("/search", handler.Search(db))

		authed.GET("/overview", handler.Overview(db))
		authed.GET("/gce", handler.GCE(db))
		authed.GET("/gae", handler.GAE(db))
		authed.GET("/sql", handler.SQL(db))
		authed.GET("/functions", handler.Functions(db))
		authed.GET("/apis_enabled", handler.APIsEnabled(db))

		authed.GET("/export", handler.Export)
		authed.GET("/authorization", handler.Authorization(db))
		authed.GET("/account", handler.Account)
		authed.POST("/logout", identity.HandleLogout)
		// ******** End left nav

		authed.GET("/revokegcp", identity.RevokeGCP(db))
		authed.GET("/disconnectgcp", identity.DisconnectGCP(db))
		authed.GET("/google/oauth2", identity.GoogleOAuth2Callback(db))

		// Fake privileged routes for demo
		authed.GET("/password", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "password.tmpl", gin.H{})
			c.Abort()
		})
		authed.GET("/encryption", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "encryption.tmpl", gin.H{})
			c.Abort()
		})
		authed.GET("/2fa", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "2fa.tmpl", gin.H{})
			c.Abort()
		})
		authed.GET("/status", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "status.tmpl", gin.H{})
			c.Abort()
		})

		authed.GET("/project.csv", func(c *gin.Context) {
			if c.Query("t") != "project" {
				c.Data(http.StatusBadRequest, "text/csv; charset=utf-8", []byte{})
				c.Abort()
				return
			}

			email, exists := c.Get(identity.IdentityContextKey)
			if !exists {
				panic("Unexpected")
			}
			var user types.User
			err := db.First(&user, "email = ?", email).Error
			if err != nil {
				panic("Unvalid DB state")
			}
			if user.AccessToken == nil {
				core.HTMLWithGlobalState(c, "overview.tmpl", gin.H{
					"Unauthorized": true,
				})
				return
			}

			var csvLines []string
			listProjectResponse, _ := gcecache.ReadProjectsCache(db, user)
			if listProjectResponse != nil {
				csvLines = append(csvLines, "project_id, project_name")
				for _, v := range listProjectResponse.Projects {
					csvLines = append(csvLines, v.ProjectId+","+v.ProjectNumber)
				}
			}
			c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(strings.Join(csvLines[:], "\n")))
			c.Abort()
		})
	}

	r.Run(":3000")
}
