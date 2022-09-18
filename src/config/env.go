package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvEnum int

const (
	Prod EnvEnum = iota
	Dev
)

var (
	Host            string
	RedirectUri     string
	EapUrl          string
	Env             EnvEnum
	CacheEnabled    bool
	Preindex        bool
	GCPClientId     string
	GCPClientSecret string
	ContactEmail    string
	BrandName       string
	CookieHttpsOnly bool
	DefaultTheme    string
	SPAEnabled      bool
)

func init() {
	Host = os.Getenv("HOST")
	if len(Host) == 0 {
		panic("Missing $HOST env variable")
	}
	fmt.Printf("Host=%v\n", Host)

	DefaultTheme = os.Getenv("DEFAULT_THEME")
	if DefaultTheme != "kubrick" {
		DefaultTheme = "requiem"
	}

	RedirectUri = os.Getenv("REDIRECT_URI")
	if len(RedirectUri) == 0 {
		panic("Missing $REDIRECT_URI env variable")
	}

	ContactEmail = os.Getenv("CONTACT_EMAIL")
	if len(ContactEmail) == 0 {
		panic("Missing $CONTACT_EMAIL env variable")
	}

	BrandName = os.Getenv("BRAND_NAME")
	if len(BrandName) == 0 {
		panic("Missing $BRAND_NAME env variable")
	}

	if os.Getenv("RELEASE_ENV") == "dev" {
		Env = Dev
		// Locally, we want to have localdev at http://
		CookieHttpsOnly = false
	} else {
		Env = Prod
		// Prod should be https://.
		CookieHttpsOnly = true
	}
	log.Printf("Env=%v\n", Env)

	v, present := os.LookupEnv("CACHE_ENABLED")
	if present {
		CacheEnabled, _ = strconv.ParseBool(v)
	}
	fmt.Printf("CacheEnabled=%v\n", CacheEnabled)

	v, present = os.LookupEnv("PREINDEX")
	if present {
		Preindex, _ = strconv.ParseBool(v)
	}
	fmt.Printf("Preindex=%v\n", Preindex)

	EapUrl = os.Getenv("EAP_URL")
	if len(EapUrl) == 0 {
		panic("Missing $EAP_URL env variable")
	}
	if !strings.HasPrefix(EapUrl, "http") {
		panic("URL must start with HTTP")
	}
	fmt.Printf("EAP_URL=%v\n", EapUrl)

	GCPClientId = os.Getenv("GCP_CLIENT_ID")
	if len(GCPClientId) == 0 {
		panic("Missing $GCP_CLIENT_ID env variable")
	}

	GCPClientSecret = os.Getenv("GCP_CLIENT_SECRET")
	if len(GCPClientSecret) == 0 {
		panic("Missing $GCP_CLIENT_SECRET env variable")
	}

	if os.Getenv("SPA_ENABLED") == "hellyes" {
		SPAEnabled = true
	} else {
		SPAEnabled = false
	}
}
