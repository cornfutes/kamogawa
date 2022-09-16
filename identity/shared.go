package identity

import "os"

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

const sessionCookieKey = "session"

const IdentityContextkey = "AuthorizedUserEmail"
