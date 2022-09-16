package config

import (
	"fmt"
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
	Host         string
	RedirectUri  string
	EapUrl       string
	Env          EnvEnum
	CacheEnabled bool
	Preindex     bool
)

func init() {
	Host = os.Getenv("HOST")
	if len(Host) == 0 {
		panic("Missing $HOST env variable")
	}
	fmt.Printf("Host=%v\n", Host)

	RedirectUri = os.Getenv("REDIRECT_URI")
	if len(RedirectUri) == 0 {
		panic("Missing $REDIRECT_URI env variable")
	}

	if len(os.Getenv("DEV")) > 0 {
		Env = Dev
	} else {
		Env = Prod
	}
	fmt.Printf("Env=%v\n", Env)

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

	fmt.Printf("EAP_URL =%v\n", CacheEnabled)

}
