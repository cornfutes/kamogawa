package config

import (
	"fmt"
	"os"
)

type EnvEnum int

const (
	Prod EnvEnum = iota
	Dev
)

var (
	Host         string
	RedirectUri  string
	Env          EnvEnum
	CacheEnabled bool
)

func init() {
	Host = os.Getenv("HOST")
	if len(Host) == 0 {
		panic("Missing $HOST env variable")
	}

	RedirectUri = os.Getenv("REDIRECT_URI")
	if len(Host) == 0 {
		panic("Missing $REDIRECT_URI env variable")
	}

	if len(os.Getenv("DEV")) > 0 {
		Env = Dev
	} else {
		Env = Prod
	}

	CacheEnabled = len(os.Getenv("CACHE_ENABLED")) > 0
	fmt.Printf("CCacheEnabled=%v\n", CacheEnabled)
}
