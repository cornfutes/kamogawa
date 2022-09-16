package main

import (
	"fmt"
	"kamogawa/asset"
	"kamogawa/cache"
	"kamogawa/core"
	"kamogawa/handler"
	"kamogawa/identity"
	"kamogawa/types"
	"log"
	"strings"
	"unsafe"

	"net/http"

	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

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
}

func main() {
	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(identity.SetAuthContext())
	asset.Config(r)

	r.GET("/ping", func(c *gin.Context) {
		q := c.Query("q")
		fmt.Printf("Query: %v\n", q)
		handler.SearchInstances(db, q)
		c.JSON(http.StatusOK, "{'result': 'pong'}")
	})

	r.GET("/login", handler.Login)
	r.GET("/reset", handler.Reset)
	r.POST("/login", identity.HandleLogin)
	r.POST("/reset", identity.HandleReset)
	r.GET("/loggedout", func(c *gin.Context) {
		if identity.ExtractClaimsEmail(c) == nil {
			core.HTMLWithGlobalState(c, "loggedout.html", gin.H{})
		} else {
			c.Redirect(http.StatusFound, "/account")
		}
	})

	authed := r.Group("/", identity.GateAuth())
	{
		authed.GET("release.txt", asset.TXT(asset.Release))

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
		authed.GET("/google/oauth2", identity.GoogleOAuth2Callback(db))

		// Fake privileged routes for demo
		authed.GET("/password", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "password.html", gin.H{})
			c.Abort()
		})
		authed.GET("/encryption", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "encryption.html", gin.H{})
			c.Abort()
		})
		authed.GET("/2fa", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "2fa.html", gin.H{})
			c.Abort()
		})
		authed.GET("/status", func(c *gin.Context) {
			core.HTMLWithGlobalState(c, "status.html", gin.H{})
			c.Abort()
		})

		authed.GET("/project.csv", func(c *gin.Context) {
			if c.Query("t") != "project" {
				c.Data(http.StatusBadRequest, "text/csv; charset=utf-8", []byte{})
				c.Abort()
				return
			}

			var email, exists = c.Get(identity.IdentityContextKey)
			if !exists {
				panic("Unexpected")
			}
			var user types.User
			err := db.First(&user, "email = ?", email).Error
			if err != nil {
				panic("Unvalid DB state")
			}
			if user.AccessToken == nil {
				core.HTMLWithGlobalState(c, "overview.html", gin.H{
					"Unauthorized": true,
				})
				return
			}

			var csvLines []string
			listProjectResponse, _ := cache.ReadProjectsCache(db, user)
			if listProjectResponse != nil {
				csvLines = append(csvLines, "project_id, project_name")
				for _, v := range listProjectResponse.Projects {
					csvLines = append(csvLines, v.ProjectId+","+v.ProjectNumber)
				}
			}
			c.Data(http.StatusOK, "text/csv; charset=utf-8", strToBytes(strings.Join(csvLines[:], "\n")))
			c.Abort()
		})
	}

	r.Run(":3000")
}

func strToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
