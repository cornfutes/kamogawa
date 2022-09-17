package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"kamogawa/cache/gcecache"
	"kamogawa/core"
	"kamogawa/handler"
	"kamogawa/identity"
	"kamogawa/media"
	"kamogawa/types"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
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
	db.FirstOrCreate(&types.User{
		Email:    "1337gamer@gmail.com",
		Password: "1234",
	})
	db.FirstOrCreate(&types.User{
		Email:    "team@otonomi.ai",
		Password: "dHJDFh43aa.X",
	})
	db.FirstOrCreate(&types.User{
		Email:    "null@hackernews.com",
		Password: "Pb$droV@a&t.a0e3",
	})
}

func main() {
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
		tokenString, err := c.Cookie(identity.SessionCookieKey)
		if err != nil {
			c.Redirect(http.StatusFound, "/account")
		}
		if identity.ExtractClaimsEmail(tokenString, c) == nil {
			core.HTMLWithGlobalState(c, "loggedout.tmpl", gin.H{})
		} else {
			c.Redirect(http.StatusFound, "/account")
		}
	})

	authed := r.Group("/", identity.GateAuth())
	{
		authed.GET("release.txt", media.Data(media.MimeTypeTXT, media.Release))

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
