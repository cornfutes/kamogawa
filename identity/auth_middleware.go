package identity

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractClaimsEmail(c *gin.Context) *string {
	tokenString, err := c.Cookie(sessionCookieKey)
	if err != nil {
		return nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecretKey, nil
	})
	if err != nil {
		log.Printf("Error validating JWT")
		return nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		emailString := fmt.Sprintf("%v", claims["email"])
		return &emailString
	} else {
		return nil
	}
}

func SetAuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := ExtractClaimsEmail(c)
		if result != nil {
			c.Set(IdentityContextKey, *result)
		}
	}
}

func GateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(IdentityContextKey)
		if !exists {
			c.Redirect(http.StatusFound, "/demo") // login?&error="+strconv.Itoa(int(Unauthorized)))
			c.Abort()
		}
	}
}
