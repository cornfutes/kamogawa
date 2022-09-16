package identity

import (
	"kamogawa/config"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type IdentityError int64

const (
	Undefined IdentityError = iota + 1
	Internal
	Incorrect
	Unauthorized
)

func (s IdentityError) String() string {
	switch s {
	case Internal:
		return "An internal error occurred."
	case Incorrect:
		return "Your email or password was incorrect."
	case Unauthorized:
		return "Invalid session. Please re-login."
	}

	return "An unexpected error occurred."
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

const testUserEmail string = "1337gamer@gmail.com"
const testUserPassword string = "HeroBallZ$5"

func HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		log.Fatal("Invalid Login Request")
	}

	// TODO: replace with DB lookup
	// TODO: encrypt password with bcrypt
	if !(loginRequest.Email == testUserEmail && loginRequest.Password == testUserPassword) {
		c.Redirect(http.StatusFound, "/login?email="+loginRequest.Email+"&error="+strconv.Itoa(int(Incorrect)))
		c.Abort()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": loginRequest.Email,
	})
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?email="+loginRequest.Email+"&error="+strconv.Itoa(int(Internal)))
		c.Abort()
		return
	}

	// TODO: change in prod. Cookie can be set over non-https domain ( i.e. http://localhost )
	// httpOnly flag set to true, preventing cookie being accessed by JavaScript
	c.SetCookie(sessionCookieKey, tokenString, 3600, "/", config.Host, false, true)

	c.Redirect(http.StatusFound, "/account")
}
