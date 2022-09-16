package handler

import "strings"

var baseScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/cloud-platform.read-only",
	"https://www.googleapis.com/auth/compute.readonly",
}

type ScopeMetadata struct {
	Scope       string
	Description string
	UsedBy      []string
}

var AllScopes = []ScopeMetadata{
	{
		Scope:       "https://www.googleapis.com/auth/userinfo.email",
		Description: "https://www.googleapis.com/auth/userinfo.email",
		UsedBy:      []string{""},
	},
	{
		Scope:       "https://www.googleapis.com/auth/cloud-platform.read-only",
		Description: "https://www.googleapis.com/auth/cloud-platform.read-only ( read-only )",
		UsedBy:      []string{"Everything"},
	},
	{
		Scope:       "https://www.googleapis.com/auth/compute.readonly",
		Description: "https://www.googleapis.com/auth/compute.readonly ( read-only )",
		UsedBy:      []string{"GCE"},
	},
	{
		Scope:       "https://www.googleapis.com/auth/cloud-platform",
		Description: "https://www.googleapis.com/auth/cloud-platform ( read/write )",
		UsedBy:      []string{"Cloud Functions", "CloudSQL"},
	},
	{
		Scope:       "https://www.googleapis.com/auth/sqlservice.admin",
		Description: "https://www.googleapis.com/auth/sqlservice.admin ( read/write )",
		UsedBy:      []string{"Cloud Functions", "CloudSQL"},
	},
}

var GCPScopes string = strings.Join(baseScopes, " ")
