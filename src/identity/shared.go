package identity

import "os"

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

const SessionCookieKey = "session"

const IdentityContextKey = "AuthorizedUserEmail"
