package identity

import (
	"log"
	"net/http"
	"strconv"

	"kamogawa/config"

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
		return "Please use the demo to login."
	}

	return "An unexpected error occurred."
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

const (
	testUserEmail     string = "1337gamer@gmail.com"
	testUserPassword  string = "HeroBallZ$5"
	testUserEmail2    string = "team@otonomi.ai"
	testUserPassword2 string = "dHJDFh43aa.X"
	testUserEmail3    string = "null@hackernews.com"
	testUserPassword3 string = "Pb$droV@a&t.a0e3"
)

func HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		log.Fatal("Invalid Login Request")
	}

	// TODO: replace with DB lookup
	// TODO: encrypt password with bcrypt
	if !((loginRequest.Email == testUserEmail && loginRequest.Password == testUserPassword) ||
		(loginRequest.Email == testUserEmail2 && loginRequest.Password == testUserPassword2) ||
		(loginRequest.Email == testUserEmail3 && loginRequest.Password == testUserPassword3)) {
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

	// httpOnly flag set to true, preventing cookie being accessed by JavaScript
	c.SetCookie(SessionCookieKey, tokenString, 3600, "/", config.Host, config.CookieHttpsOnly, true)

	c.Redirect(http.StatusFound, "/search?q=test")
}
